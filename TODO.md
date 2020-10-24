create and use base response object model
context between processes: context extract headers
amqp transports
context timeout and object (now I return max, need to change it)
// TODO: we need to calculate the deadline and timeout for the callee, so there should be some substruction
functions in core.url should be moved to the encoding_decoding

Rabbit Payload should be similar to:
{
  "messageId": "b7980000-56be-0050-2820-08d866d872a6",
  "correlationId": "260c3944-0997-4afc-93a3-13b8e94e30e4",
  "conversationId": "b7980000-56be-0050-2a54-08d866d872a6",
  "sourceAddress": "rabbitmq://tle-rabbitmq-ha-chart-0.tle-rabbitmq-ha-chart-discovery.default.svc.cluster.local:32672/thelotter/bus-INT-WIN3-TheLotter.Marketing.ConversionRetention.Index.Service-s6cyyynszayfb43xbdcgpuwi8w?durable=false&autodelete=true",
  "destinationAddress": "rabbitmq://tle-rabbitmq-ha-chart-0.tle-rabbitmq-ha-chart-discovery.default.svc.cluster.local:32672/thelotter/TheLotter.Marketing.ConversionRetention.Index.Shared:RunETLProcessCommand",
  "messageType": [
    "urn:message:TheLotter.Marketing.ConversionRetention.Index.Shared:RunETLProcessCommand",
    "urn:message:TheLotter.Core2.Messaging:CommandMessage",
    "urn:message:TheLotter.Core2.Messaging:MessageBase",
    "urn:message:TheLotter.Core2.Messaging:IMessage",
    "urn:message:TheLotter.Core2:INonCachableKey",
    "urn:message:TheLotter.Core2:IMessageCallContextInfo",
    "urn:message:TheLotter.Core2.Messaging:ICommandMessage",
    "urn:message:TheLotter.Marketing.ConversionRetention.Index.Shared:IRunETLProcessCommand"
  ],
  "message": {
    "request": {
      "dateFrom": "2020-01-10T00:00:00",
      "dateTo": "2020-03-10T00:00:00",
      "etlTypes": [
        2
      ]
    },
    "senderTimeoutInMiliseconds": 30000.0,
    "timestamp": "2020-10-02T13:38:29.6339924Z",
    "correlationId": "260c3944-0997-4afc-93a3-13b8e94e30e4"
  },
  "headers": {
    "use-local-cache": {
      "data": true
    }
  },
  "host": {
    "machineName": "INT-WIN3",
    "processName": "TheLotter.Marketing.ConversionRetention.Index.Service",
    "processId": 70316,
    "assembly": "TheLotter.Marketing.ConversionRetention.Index.Service",
    "assemblyVersion": "1.0.2.191",
    "frameworkVersion": "4.0.30319.42000",
    "massTransitVersion": "3.5.7.1082",
    "operatingSystemVersion": "Microsoft Windows NT 6.2.9200.0"
  }
}