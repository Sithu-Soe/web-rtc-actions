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

//HTML elements
const webcamButton = document.getElementById('webcamButton');
const webcamVideo = document.getElementById('webcamVideo');
const callButton = document.getElementById('callButton');
const callInput = document.getElementById('callInput');
const answerButton = document.getElementById('answerButton');
const remoteVideo = document.getElementById('remoteVideo');
const hangupButton = document.getElementById('hangupButton');


webcamButton.onclick = async () => {
  console.log("HEHE");

  localStream = await navigator.mediaDevices.getUserMedia({ video: true, audio: true})

  localStream.getTracks().forEach((track) => {
    pc.addTrack(track, localStream)
  })

  // remoteStream = new MediaStream()

  webcamVideo.srcObject = localStream;

  callButton.disabled = false;
  answerButton.disabled = false;
  webcamButton.disabled = true;
}