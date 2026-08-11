

#规则转发
## 连接状态规则转发
```
SELECT
  clientid,
  event,
  timestamp,
  username as device_key,
  reason
FROM
  "$events/client/connected",
  "$events/client/disconnected"
```

### connect
```
{"node":"emqx@172.17.0.3","timestamp":1769667983517,"peername":"127.0.0.1:52918","sockname":"0.0.0.0:1883","keepalive":60,"metadata":{"namespace":"global","rule_id":"sql_tester:b4c542c5886fd367"},"event":"client.connected","username":"u_emqx","proto_ver":5,"client_attrs":{"test":"example"},"mountpoint":"undefined","clientid":"c_emqx","connected_at":1769667983517,"is_bridge":false,"proto_name":"MQTT","clean_start":true,"expiry_interval":3600,"conn_props":{"User-Property":{"foo":"bar"},"User-Property-Pairs":[{"key":"foo"},{"value":"bar"}],"Session-Expiry-Interval":7200,"Receive-Maximum":32}}

```
### disconnect
```
{"node":"emqx@172.17.0.3","reason":"normal","timestamp":1769667859192,"peername":"192.168.0.10:56431","sockname":"0.0.0.0:1883","metadata":{"namespace":"global","rule_id":"sql_tester:5101a978f8e3cee2"},"event":"client.disconnected","username":"u_emqx","proto_ver":5,"client_attrs":{"test":"example"},"clientid":"c_emqx","connected_at":1769667859192,"proto_name":"MQTT","disconn_props":{"User-Property":{"foo":"bar"},"User-Property-Pairs":[{"key":"foo"},{"value":"bar"}],"Session-Expiry-Interval":7200,"Reason-String":"Redirect to another server","Server Reference":"192.168.22.129"},"disconnected_at":1769667859192}
```