# CFTools Relay

[![Discord](https://img.shields.io/discord/729467994832371813?color=7289da&label=Discord&logo=discord&logoColor=ffffff&style=flat-square)](https://go2tech.de/discord)
[![Tests](https://img.shields.io/github/workflow/status/FlorianSW/cftools-relay/build?label=tests&style=flat-square)](https://github.com/FlorianSW/cftools-relay/actions/workflows/build.yml)

CFTools Relay is an easy-to-use, still in development, tool that allows you to subscribe to CFTools Cloud Webhook events and forward them to a different target.
Right now, the only target that is supported is as Discord Webhook URL.

## Why?

CFTools Relay is a tool that is mostly done to allow filtering for specific CFTools events.
CFTools Cloud already provides a way to send Webhooks for specific events to Discord, however, apart from the event name, there is no additional filter criteria that can be applied.

If one, e.g., would like to get a notification in Discord only, if a player kills another player from more than 500 meters away, potentially with a specific weapon (like an IJ 70) only, this is currently not possible with CFTools Webhooks.
CFTools Relay allows to define such filters, which are applied to the incoming Webhook events.
Only if at least one filter, with all filter rules, match the event, it will be forwarded/relayed to Discord.

## Should I use this tool?

If you do not need to filter Webhook events based on data that is within an event (rather than just the event type), then you SHOULD NOT use this tool.
It is only increasing complexity in this case, and you can achieve the very same outcome with the CFTools Cloud builtin Webhook-to-Discord feature.
Simply use that.

If you, however, need some more filter logic on top of your Webhook events to further reduce the amount of messages your channel is flooded with, this tool might be of help.

## Installation

The installation is as simple as downloading the latest version of this tool, put it somewhere and start it, either by double-clicking on the exe file (for Windows OS) or running the tool in a terminal (for Linux and Windows OS).

### Download the latest version

You can either visit the [release page](https://github.com/FlorianSW/cftools-relay/releases) to download the latest released version (which might miss some features from the current development version).

Alternatively, you can visit the [automated build page](https://github.com/FlorianSW/cftools-relay/actions/workflows/build.yml), select the most recent successful run (with a green checkmark) and download the artifact for your operating system from the Artifacts section.

## Configuration

Upon the first start of the tool, it will automatically create a configuration file, named `config.json` in the same directory as the binary.

### Check port of the tool and your firewall

This tool will start an internal webserver, which is used to handle webhook events sent from CFTools to it.
The webserver, by default, uses port 8080 to listen on.
You need to ensure that this port is available to the tool and not used by another program/process already (like another webserver).
If you do not want to use port 8080 for this tool, you can change the port in the created configuration file.

Also, make sure that the port you configured (or the default one) is whitelisted in your firewall configuration.
The CFTools Relay is using an unencrypted connection using http.
If you want to use a TLS-encrypted connection for your webhook messages, you may want to setup a reverse proxy for this tool, which handles the TLS termination.
A setup like that is out-of-scope of this README.

### First time setup

When you setup this tool the first time, you need to do some basic steps in order to create and verify the webhook in the CFTools Cloud console:

1. Start the tool (if you did not do that already) and ensure it is able to receive web requests (see _Check port of the tool and your firewall_)
2. Go to the CFTools Cloud Dashboard and open the server where the webhook should be added to
3. Navigate to the Manage -> Integrations page of that server
4. Add a new Webhook with the _New Webhook_ button
5. Enter the Webhook URL:
   - It consists of the public IP address or domain of your server (e.g. http://123.123.123.123)
   - appended is the Port (by default 8080) separated by a colon (:8080)
   - at the end it needs to have a fixed path: `/cftools-webhook`
   - For the example values above, the full webhook URL looks like: `http://123.123.123.123:8080/cftools-webhook`
6. Select `CFTools CLoud (Hephaistos v1)` as the _Payload format_
7. Click _Deploy_
8. Reload the page and make sure it shows a green shield (which means _Verified and active_) next to the newly created webhook
9. Open the Webhook details and copy the value in the _Secret_ field
10. Open the `config.json` of CFTools Relay in your favourite text editor
11. Paste the copied secret into the value of the `secret` field inside the config. It should then look like this (when your secret is `abc123`):
    - `"secret": "abc123",`
12. Save the `config.json` file and restart the CFTools Relay tool
13. Go back to the Webhook details page of CFTools Cloud and select every event you want the CFTools Relay to receive
14. Hit _Save_

You're done with the configuration on the CFTools side.
Given CFTools Relay does not know where to relay the webhook event messages to right now, you need to configure the Discord target:

1. Go to the Discord channel of your choice where the Webhook messages should be relayed to
2. Open the settings of this channel
3. Navigate to the _Integrations_ section of the settings
4. Create a new webhook (or use an existing one, up to you)
5. Copy the _Webhook URL_ value
6. Open the `config.json` of CFTools Relay in your favourite text editor 7Paste the copied URL into the value of the `discord.webhook_url` field inside the config. It should then look like this (when your URL is `http://example.com`):
```json
"discord": {
  "webhook_url": "http://example.com"
}
```
7. Save the `config.json` file and restart the CFTools Relay tool

That's it.
The first time configuration is now done and CFTools Relay will now start to forward Webhook events from CFTools Cloud to your discord channel.

### Filter configuration

The main use case of CFTools Relay is to be able to filter events that should be forwarded to Discord.
By default, there are no filters set up, which makes CFTools Relay to relay all messages to Discord.

There are two main concepts to understand before using filters:

**Filters:**
You can have 0-n filters configured in CFTools Relay.
Filters are evaluated independently of each other.
If at least one filter matches a specific webhook event received from CFTools Cloud, this event will be relayed to Discord.
That means, filters are combined with an OR operator.
The only required field for filters is the event name this filter should apply to.

**Rules:**
Each filter can have 0-n rules that are evaluated to the webhook event.
Only if all defined rules match the event, the filter evaluates to "match the event", which would make the event be relayed to Discord.
Rules are therefore combined with an AND operator.

Rules allow you to define more in-depth filtering on Webhook events, as they allow you to compare specific fields to values you define.
The available fields heavily depend on the Webhook event type (see the JSON files in the `payloads/` folder of this project for some information).

A usual configuration with one example filter looks like:

```json
{
  "port": 8080,
  "secret": "...",
  "discord": {
    "webhook_url": "..."
  },
  "filter": [
    {
      "event": "user.join",
      "rules": null
    },
    {
      "event": "user.leave",
      "rules": null
    }
  ]
}
```

This filter will make CFTools Relay to only relay events with the type `user.join` and `user.leave`.

#### Example 1: Relay kill event when distance is greater than 1000 meter

If you wish to get a kill notification in your Discord channel only, if the distance between the _Victim_ and the _Murderer_ is greater than 1000 meters, you can add the following filter:

```json
{
  "event": "player.kill",
  "rules": [
    {
      "comparator": "gt",
      "field": "distance",
      "value": 1000
    }
  ]
}
```

#### Example 2: Relay kill event when distance is greater than 100 meter with an IJ-70

Based on the previous example, you can also combine multiple rules and make filters that, e.g., relay a message to Discord only, if the kill was made over 100 meters between the Victim and the Murderer, as well as if the murderer used an IJ-70:

```json
{
  "event": "player.kill",
  "rules": [
    {
      "comparator": "gt",
      "field": "distance",
      "value": 100
    },
    {
      "comparator": "eq",
      "field": "weapon",
      "value": "IJ-70"
    }
  ]
}
```

#### Available Comparators

The following table lists the available Comparators for filter rules:

| Comparator   | Explanation |
|--------------|-------------|
| `eq`         | `Equals` comparator, matches only, if the value of the configured `field` in the event is exactly the configured `value`. This comparator is case-sensitive. |
| `gt`         | `Greater than or equals` comparator, matches only, if the value of the configured `field` in the event is a numeric value and is greater than or equals the configured `value`. This comparator never matches when the value is not a numeric field. |
| `lt`         | `Less than or equals` comparator, matches only, if the value of the configured `field` in the event is a numeric value and is les than or equals the configured `value`. This comparator never matches when the value is not a numeric field. |
| `contains`   | `Contains` comparator, matches only, if the value of the configured `field` in the event is a contains the configured `value`. This is a wildcard matcher, which is equivalent to `*value*` is wildcards would exist. |
| `startsWith` | `Starts with` comparator, same as `contains` with the only difference, that the field value needs to start with the configured value (`value*` instead of `*value*`). |
| `endsWith`   | `Ends with` comparator, same as `contains` with the only difference, that the field value needs to end with the configured value (`*value` instead of `*value*`). |

## Virtual Fields

In addition to the fields received via the webhook from CFTools, CFTools Relay may inject available data for the specific event in it's own fields.
These fields are also called _virtual fields_ as they are only present and usable within CFTools Relay rule configuration.

Virtual Fields aim to provide context information for an event within the stream of incoming webhook events.
They allow to have filters match when, e.g., a specific threshold of events is reached, instead of relaying each and every webhook event.

You can use each of the below listed virtual fields as the `field` name in any rule configuration.
As virtual fields may be related to the time-series of events, you need to specify the timeframe for which the virtual field should be evaluated on.
If you do not specify this value, a default of 1 hour is assumed.
The time-frame is the amount of time and a unit, combined in a single string, e.g. "1h" for one hour or "30m" for 30 minutes.
Available units can be found in [this documentation](https://pkg.go.dev/time#ParseDuration).

Events are associated to the CFTools ID they are about and are only recognized for virtual fields of the same player.
For events that may have multiple CFTools IDs (like `player.kill`, which has a victim and a murderer CFTools ID), the event is persisted for the CFTools ID, which matches the event the most (`player.kill` is an event of a kill, which is associated with the murderer).
For each CFTools ID, a maximum of 100 events are preserved before the most historical events are truncated/removed.

### Available Virtual Fields

| Field name          | Explanation |
|---------------------|-------------|
| `vf_event_count`    | A simple counter. Counts every event with the same `type` (e.g. `player.kill` for the specified time-frame. |

#### Example 1: Relay kill message only, when 5 kills of the same player within the last hour

The following example will only relay the Webhook to Discord, if the murderer had 5 kills in the last hour (including the currently evaluated one):

```json
{
  "event": "player.kill",
  "rules": [
    {
      "comparator": "gt",
      "field": "vf_event_count",
      "value": 5,
      "since": "1h"
    }
  ]
}
```

#### Example 2: Relay 5th kill message only

Based on the previous example, the set of rules in this example will _only_ relay the 5th kill message.
Kill events that appear after that are ignored.

```json
{
  "event": "player.kill",
  "rules": [
    {
      "comparator": "gt",
      "field": "vf_event_count",
      "value": 5,
      "since": "1h"
    },
    {
      "comparator": "lt",
      "field": "vf_event_count",
      "value": 5,
      "since": "1h"
    }
  ]
}
```

## Custom message & color

When relaying a message to your discord, CFTools Relay uses a default message, which depends on the type of event.
It also chooses a default color for the embedded message, which is currently dark blue.

Optionally, you can define a custom message as well as a custom color for the message on a filter, which is used whenever the filter matches.
Both options are optional and do not depend on each other.
When defining a custom message for a filter, all the metadata of the message will still be relayed to your discord channel.

### Example 1: Configure a custom message for a filter

The following filter configuration will use `A custom message` when the first filter matches, and the default when the second matches.

```json
  "filter": [
    {
      "event": "user.join",
      "rules": null,
      "message": "A custom message"
    },
    {
      "event": "user.leave",
      "rules": null
    }
  ]
```

### Example 1: Configure a custom color for a filter

The following filter configuration will use `A custom message` when the first filter matches, and the default when the second matches.

```json
  "filter": [
    {
      "event": "user.join",
      "rules": null,
      "color": "DARK_GREEN"
    },
    {
      "event": "user.leave",
      "rules": null
    }
  ]
```

### Supported color names

You can choose from a pre-configured color palette.
Supported color names can be found in the [source code](./internal/domain/filter.go).
