const express = require('express');
const res = require('express/lib/response');
const axios = require('axios') 
const app = express();
const port = 3000;


// Get Method
app.get('/', (req, res) => {
    res.send('Hello World!');
});
app.get('/db', (req, res) => {
    axios.get("https://localhost:8080/chatbot")
        .then(outcomes => {
            res.json(outcomes.data) 
        })
        .catch(err => {
            console.log(err);
            reportError();
        })
});

// Post Method
app.post("/db", (req, res) => {
    var id = req.body.id;
    var title = req.body.title;
    var solution = req.body.solution;
    var count = req.body.count;
    var additional = req.body.additional;   
    var count = req.body.count1;


    var data = new collections({
        "id": id,
        "title": title,
        "solution": solution,
        "count": count,
        "additional": additional,
        "count": count
    });
    data.save(function (err, results) {
        if (err) {
            console.log(err);
        } else {
            console.log(results);
        }
    });
});

app.listen(port, () => {
    console.log(`Server is running on  ${port}`)
});