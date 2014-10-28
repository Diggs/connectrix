## Connectrix

Connectrix is a general purpose tool for sending and receiving events to and from computer systems. Think of it like IFTTT for DevOps.

### Using Connectrix

Connectrix ships with a series of general purpose event sources and sinks (known as ```Channels```) that let it send and receieve events. Currently the supported mechanisms are HTTP(S) and IRC. Things like SNMP and SMTP are planned as well.

Connectrix works at a fairly low level (e.g. dealing with HTTP) and allows you to build on top of that by defining different event sources and types in a declaritve way. Connectrix will gather raw info about the incoming event and then evaluate each of the sources and types that have been delcared until it finds a match. For example you may define an event source as ```GitHub``` and an event type of ```push```.

### Defining event sources

Event sources are declared in the config.json file, the options are:

 * name - the name of the event source
 * hint - a string to match against the raw event info to identify the event source (see docs for each channel to see what the hints are that can be matched against)
 * parser - the name of the parser that should be used to parse the event data (json, xml and yaml are supported)
 * events - a list of events that the source will send (see next section)

Here's an example of using GitHub as an event source. Github sends an HTTP User-Agent header of 'GitHub-Hookshot' so that can be used to identify it. GitHub sends JSON data in the HTTP body so we tell Connectrix to use the JSON parser.

```
"sources":[
	{
		"name":"GitHub",
		"hint":"User-Agent:GitHub-Hookshot",
		"parser":"json",
		"events":[]
	}
]
```

If the system you're defining these rules for doesn't send anything that uniquely identifies it by default then perhaps you can embed an indentifier in the URL it's sending data to. Here's an example of doing that with ```CircleCI```:

```
"sources":[
	{
		"name":"CircleCI",
		"hint":"source=circleci",
		"parser":"json",
		"events":[]
	}
]
```

Now when you set up the webook in CircleCI you just add ?source=circleci as a query parameter and Connectrix will be able to identify that it came from CircleCI.

### Defining event types

Once an event source has been declared you can then specify each of the event types that source will send, the options for event types are:

 * type - the name of the event type
 * hint - a string to match against the raw event info to identify the event type (see docs for each channel to see what the hints are that can be matched against)
 * template - a [go template](http://gohugo.io/templates/go-templates/) compatible string that the event content will be run through to generate a human readable representation of the event

Here's an extended GitHub example with the push event type declared. Notice again that GitHub sets the X-Github-Event HTTP header that we can use to identify the event type.

```
"sources":[
	{
		"name":"GitHub",
		"hint":"User-Agent:GitHub-Hookshot",
		"parser":"json",
		"events":[
			{
				"type":"push",
				"hint":"X-Github-Event:push",
				"template":"{{.pusher.name}} committed to {{.repository.name}}:{{.ref}}: {{.head_commit.message}} - {{.head_commit.url}}"
			}
		]
	}
]
```

Here's the extended CircleCI example, again making use of an additional query param to identify the event type:

```
"sources":[
	{
		"name":"CircleCI",
		"hint":"source=circleci",
		"parser":"json",
		"events":[
			{
				"type":"build",
				"hint":"event=build",
				"template":"Build #{{.payload.build_num}} of {{.payload.reponame}}:{{.payload.branch}} by {{.payload.committer_name}} finished with status: {{.payload.outcome}} - {{.payload.build_url}}"
			}
		]
	}
]
```

### Routing events

Once event sources and types have been declared routes can be defined that tell Connectrix what to do when it recieves an event. Typically you would route the event from one Channel to another. For example you might say "if a build fails in CircleCI open a Github issue":

```
"routes":[
	{
		"namespace":"0",
		"event_source":"CircleCI",
		"event_type":"build",
		"sub_channel_name":"http",
		"sub_channel_args":{"URL":"https://api.github.com/repos/{{.payload.username}}/{{.payload.reponame}}/issues", "Headers":"Authorization:Basic base64(username:pass)"},
		"template":"{\"title\":\"Build Failed\", \"body\":\"\"}",
		"rule":"`{{.payload.status}}` == `failed`"
	}
]
```

or "when someone commits to my GitHub repo annouce it in IRC":

```
"routes":[
	{
		"namespace":"0",
		"event_source":"GitHub",
		"event_type":"push",
		"sub_channel_name":"irc",
		"sub_channel_args":{"IRC Server":"irc.freenode.net", "IRC Channel":"#connectrix", "Nickname":"connectrix-bot"},
		"template":"",
		"rule":"`{{.repository.name}}` == `connectrix`"
	}
]
```

The options for routes are:

 * namespace - ignore this for now, always set it to 0 (it will be used to support multi-tenanted Connextrix installs in future)
 * event_source - the name of the event source you want to route
 * event_type - the type of event (from the specified event source) you want to route
 * sub_channel_name - the name of the channel to route throug (e.g. http or irc)
 * sub_channel_args - arguments to pass to the channel when routing (see docs for each channel to see what args they accept)
 * template - an optional template to run the event data through before sending it to the channel. Leave blank to use the default template specified on the event type.
 * rule - a binary expression to decide if the event should be routed (this is templated prior to being evaluated, so you can use the event data to determine if the event should be routed)

### HTTP Channel

The HTTP channels allows events to be sent and received over HTTP(S).

#### Hints

The HTTP channel uses all HTTP headers and query paramters as hints. Headers are seperated by colons and query params by equal signs.

For example if an HTTP request is received with the headers:

```
User-Agent GitHub
Content-Type application/json
```

Then you could match it with the hint "User-Agent:GitHub" or "Content-Type:application/json".

If an HTTP request is received with the URL:

```
http://foo.com/bar?source=circle&event=build
```

Then you could match it with the hint "source=circleci" or "event=build"

#### Args

 * URL - The URL to post to.
 * Headers - A comma seperated listed of headers as HeaderName:HeaderValue
 * Self Signed Cert - Set to true to allow URL to be using a self signed cert

### IRC Channel

#### Hints

#### Args
