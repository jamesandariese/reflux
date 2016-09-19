package reflux

import (
	"flag"
	"encoding/json"
	"fmt"
	ic "github.com/influxdata/influxdb/client/v2"
	"net/url"
	"os"
	"time"
)

func parseInfluxUrl(influxUrl string) (scheme, host, user, password, database string, reterr error) {
	u, err := url.Parse(influxUrl)
	if err != nil {
		reterr = err
		return
	}

	user = ""
	password = ""
	pwset := false

	if u.User != nil {
		user = u.User.Username()
		password, pwset = u.User.Password()
	}
	if user == "" {
		user = os.Getenv("INFLUX_USER")
	}
	if !pwset {
		password = os.Getenv("INFLUX_PWD")
	}
	scheme = u.Scheme
	host = u.Host
	if len(u.Path) == 0 {
		database = ""
	} else {
		if u.Path[0] == '/' {
			database = u.Path[1:]
		} else {
			database = u.Path
		}
	}
	return
}

type Client struct {
	client   ic.Client
	database string
	tags     map[string]string
	bp       ic.BatchPoints
}

// Reset all queued points and tags
func (c *Client) Reset() error {
	err := c.resetBatchPoints()
	c.tags = map[string]string{}

	return err
}

func (c *Client) resetBatchPoints() error {
	bp, err := ic.NewBatchPoints(ic.BatchPointsConfig{
		Database:  c.database,
		Precision: "s",
	})
	c.bp = bp
	return err
}

// Create a new client and create the database in influx
func NewClient(influxUrl string) (*Client, error) {
	scheme, host, user, password, database, err := parseInfluxUrl(influxUrl)
	if err != nil {
		return nil, err
	}
	ret := &Client{
		database: database,
	}

	httpConfig := ic.HTTPConfig{
		Addr:      fmt.Sprintf("%s://%s", scheme, host),
		Username:  user,
		Password:  password,
		UserAgent: "reflux",
	}

	if client, err := ic.NewHTTPClient(httpConfig); err != nil {
		return nil, err
	} else {
		ret.client = client
	}

	q := ic.NewQuery(`CREATE DATABASE "`+database+`"`, "", "")
	if response, err := ret.client.Query(q); err != nil {
		ret.client.Close()
		return nil, err
	} else if response.Error() != nil {
		ret.client.Close()
		return nil, response.Error()
	}
	if err := ret.Reset(); err != nil {
		ret.client.Close()
		return nil, err
	}

	return ret, nil
}

// Set tags to be used in the next call to AddPoint
func (c *Client) SetTags(tags map[string]string) {
	c.tags = tags
}

// Set tags to be used in the next call to AddPoint using a JSON dictionary
func (c *Client) SetTagsJson(jsonTags string) error {
	err := json.Unmarshal([]byte(jsonTags), &c.tags)
	return err
}

// Add a point to the next batch to be flushed
func (c *Client) AddPoint(name string, fields map[string]interface{}) error {
	pt, err := ic.NewPoint(name, c.tags, fields, time.Now())
	if err != nil {
		return err
	}
	c.bp.AddPoint(pt)
	return nil
}

// Flush all batched points.  Batched points are not automatically flushed.
func (c *Client) Flush() error {
	if err := c.client.Write(c.bp); err != nil {
		return err
	}
	if err := c.resetBatchPoints(); err != nil {
		return err
	}
	return nil
}

// Close the underlying influx client
func (c *Client) Close() error {
	return c.client.Close()
}

// Convenience function which calls NewClient, SetTagsJson, AddPoint, Flush, and Close.
func SendPointWithJsonTags(influxUrl, name string, fields map[string]interface{}, jsonTags string) error {
	c, err := NewClient(influxUrl)
	if err != nil {
		return err
	}
	defer c.Close()

	if err := c.SetTagsJson(jsonTags); err != nil {
		return err
	}
	if err := c.AddPoint(name, fields); err != nil {
		return err
	}
	if err := c.Flush(); err != nil {
		return err
	}
	return nil
}

var usingFlags = false
var influxUrlFlag string
var influxJsonTagsFlag string

// Prepare flags for flag package.  This must be called before flag.Parse.
// Sets up -influx-url with default of http://localhost:8086/{database}
// Sets up -influx-json-tags with default of {}
func PrepareFlags(database string) {
	if flag.Parsed() {
		panic("must call reflux.PrepareFlags before flag.Parse")
	}
	usingFlags = true
	flag.StringVar(&influxUrlFlag, "influx-url", "http://localhost:8086/"+database, "URL for influx")
	flag.StringVar(&influxJsonTagsFlag, "influx-json-tags", "{}", "JSON formatted influx tags")
}

// Convenience function used with PrepareFlags and SendPointWithJsonTags.
func SendPointUsingFlags(name string, fields map[string]interface{}) error {
	if !usingFlags {
		panic("must call reflux.PrepareFlags before calling SendPointFromFlags")
	}
	
	return SendPointWithJsonTags(influxUrlFlag, name, fields, influxJsonTagsFlag)
}
