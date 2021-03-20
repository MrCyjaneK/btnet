// IDEA: Allow user to preconfigure which trackers to use in btnet binary.
var trackers = [
    "wss://tracker.openwebtorrent.com"
]
var myid = makeid(16)

var peers = []

function trackermsg(msg) {
    data = JSON.parse(msg.data)
    if (!data.offer) {
        return
    }
    console.log("OFFER", data.offer)
    if (data.offer.type === "offer") {
        pid = peers.length
        peers[pid] = new SimplePeer()
        peers[pid].onmessage = console.log
        peers[pid].onerror = ((err) => {
            console.error("p2", err)
        })
        peers[pid].on('data', data => {
            console.log('got a message from peer: ' + data)
        })
        peers[pid].on('connect', () => {
            console.log("Connected!")
            // wait for 'connect' event before using the data channel
            peers[pid].send('hey, how is it going? I\'m '+myid)
        })
        peers[pid].signal(data.offer)
        console.log("connecting... peers["+pid+"]")
    }
}

function makeid(length) {
    var result           = '';
    var characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for ( var i = 0; i < length; i++ ) {
       result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
 }

function announce(anndata, tracker) {
    ws = new WebSocket(tracker)
    //ws.onopen = console.log
    ws.onerror = console.error
    ws.onmessage = trackermsg
    ws.onopen = (() => {
        ws.send(JSON.stringify(anndata))
    })
}

const p = new SimplePeer({
    initiator: true,
    trickle: false,
    config: { iceServers: 
        [
            { urls: 'stun:stun1.l.google.com:19302' },
            { urls: 'stun:stun1.l.google.com:19305' },
            { urls: 'stun:stun2.l.google.com:19302' },
            { urls: 'stun:stun2.l.google.com:19305' },
            { urls: 'stun:stun3.l.google.com:19302' },
            { urls: 'stun:stun3.l.google.com:19305' },
            { urls: 'stun:stun4.l.google.com:19302' },
            { urls: 'stun:stun4.l.google.com:19305' }
        ]
    } 
})

p.on('error', err => {
    console.error('p', err)
})

p.on('signal', data => {
    console.log('SIGNAL', data)
    if (data.type === 'offer') {
        trackeranndata = {
            "numwant": 10, // .offers.length
            "uploaded": 0,
            "downloaded": 0,
            "left": null,
            "action": "announce",
            "info_hash": "c293234b812651ef53919accd5b9a629ae155236", // Not quite sure how it is
            "peer_id": "-WW0014-iT6/"+makeid(8),
            "offers": [
                {
                    "offer": data,
                    "offer_id": makeid(8) // I'm not quite sute what is this
                }
            ]
        }
        // Now that we know what to send to tracker let's do that!
        announce(trackeranndata, trackers[0])
    }
})

p.on('data', data => {
    // got a data channel message
    console.log('got a message from peer: ' + data)
})