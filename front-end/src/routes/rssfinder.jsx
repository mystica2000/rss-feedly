import { useState } from "react";

function RSSFinder() {

    const [url,setURL] = useState("");

    const [results,setResults] = useState("")


    function handleSubmit() {

        fetch("http://localhost:8080/rss-finder",{
            method:'POST',
            headers:{"Content-Type":"text/plain"},

            body: url
        }).then((res)=> {
            res.json().then(result => {

                console.log("RSS FINDER!",result)
                setResults(result)
            })
        })

    }

    return (
       <div>
        <input type="text" onChange={(e)=>(setURL(e.target.value))}   />
        <button onClick={handleSubmit}>Search</button>
        <span>{results}</span>
       </div>
    )
}

export default RSSFinder;