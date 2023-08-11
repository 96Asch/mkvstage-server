const express = require('express')
const app = express()

app.listen(9080)

// respond with "hello world" when a GET request is made to the homepage
app.get('/', (req, res) => {
  res.send('hello world')
})
