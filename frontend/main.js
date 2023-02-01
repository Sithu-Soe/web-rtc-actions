import './style.css';

const iceConfiguration = {
  iceServers: [
    {
      urls: ['stun:stun1.l.google.com:19302', 'stun:stun2.l.google.com:19302'],
    },
  ],
  iceCandidatePoolSize: 10,
};

// Global State

const pc = new RTCPeerConnection(iceConfiguration);
let localStream = null;
let remoteStream = null;

let socket = new WebSocket("ws://localhost:8083/ws");
console.log("Attemtpting websocker connection...")

//HTML elements
const webcamButton = document.getElementById('webcamButton');
const webcamVideo = document.getElementById('webcamVideo');
const callButton = document.getElementById('callButton');
const callInput = document.getElementById('callInput');
const answerButton = document.getElementById('answerButton');
const remoteVideo = document.getElementById('remoteVideo');
const hangupButton = document.getElementById('hangupButton');

// 1. Setup media sources

webcamButton.onclick = async () => {
  console.log("HEHE");

  localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true })

  // Push tracks from local stream to peer connection
  localStream.getTracks().forEach((track) => {
    pc.addTrack(track, localStream)
  })

  // remoteStream = new MediaStream()

  webcamVideo.srcObject = localStream;

  callButton.disabled = false;
  answerButton.disabled = false;
  webcamButton.disabled = true;
}

callButton.onclick = async () => {

  let iceCandidate = null
  pc.onicecandidate = (event) => {
    if (event.candidate) {
      console.log(event.candidate, "DHDS", event.candidate.toJSON);
      iceCandidate = event.candidate
    }
  };

  const offerDescription = await pc.createOffer();
  await pc.setLocalDescription(offerDescription)

  console.log(offerDescription, "offer description")

  socket.send(offerDescription)
}