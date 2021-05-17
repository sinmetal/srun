# srun
Spanner へのアクセスを適当に試すやつ

```
fetch("/tweetUpdateDML",
{
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
    method: "POST",
    body: JSON.stringify({"ids":["00000001-02a8-43cc-bd2b-933b9dfa795b", "00000002-1c18-4bf9-bc51-56c7350ff7cc", "00000002-d3b1-4565-855f-5692e63ef131", "00000005-5c7e-4174-a314-f2e97259cd97"], "content":"tweetUpdateDML" })
})
.then(response => console.log(response))
.then(data => {
  console.log('Success:', data);
})
.catch(function(res){ console.log(res) })
```

```
fetch("/tweetUpdateBatchDML",
{
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
    method: "POST",
    body: JSON.stringify({"ids":["00000001-02a8-43cc-bd2b-933b9dfa795b", "00000002-1c18-4bf9-bc51-56c7350ff7cc", "00000002-d3b1-4565-855f-5692e63ef131", "00000005-5c7e-4174-a314-f2e97259cd97"], "content":"tweetUpdate" })
})
.then(response => console.log(response))
.then(data => {
  console.log('Success:', data);
})
.catch(function(res){ console.log(res) })
```

```
fetch("/tweetUpdate",
{
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
    method: "POST",
    body: JSON.stringify({"ids":["00000001-02a8-43cc-bd2b-933b9dfa795b", "00000002-1c18-4bf9-bc51-56c7350ff7cc", "00000002-d3b1-4565-855f-5692e63ef131", "00000005-5c7e-4174-a314-f2e97259cd97"], "content":"tweetUpdate" })
})
.then(response => console.log(response))
.then(data => {
  console.log('Success:', data);
})
.catch(function(res){ console.log(res) })
```

```
fetch("/tweetUpdateAndSelect",
{
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
    method: "POST",
    body: JSON.stringify({"id":"00000001-02a8-43cc-bd2b-933b9dfa795b" })
})
.then(response => response.json())
.then(data => {
  console.log('Success:', data);
})
.catch(function(res){ console.log(res) })
```


```
fetch("/tweetUpdateDMLAndSelect",
{
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json',
    },
    method: "POST",
    body: JSON.stringify({"id":"00000001-02a8-43cc-bd2b-933b9dfa795b" })
})
.then(response => response.json())
.then(data => {
  console.log('Success:', data);
})
.catch(function(res){ console.log(res) })
```
