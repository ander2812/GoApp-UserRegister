const express = require('express')



const app = express()

app.set('view engine', 'pug',  )

app.get("/",(req, res) => {
    res.sendFile("index.html",{root: __dirname+"/views/"})
})

app.get('/singup', (req,res) =>{
    res.sendFile("singup.html",{root: __dirname + "/views/"})
} )


app.listen(8080)