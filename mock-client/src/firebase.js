// Import the functions you need from the SDKs you need
import { initializeApp } from "firebase/app";
import { getAuth, signInAnonymously, GoogleAuthProvider, signInWithPopup } from "firebase/auth";
// TODO: Add SDKs for Firebase products that you want to use
// https://firebase.google.com/docs/web/setup#available-libraries

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
const firebaseConfig = {
    apiKey: "AIzaSyClRQZPhQa0DhmL8-vbqbrUGCY_4qHB7-s",
    authDomain: "rr-auth-467b4.firebaseapp.com",
    projectId: "rr-auth-467b4",
    storageBucket: "rr-auth-467b4.appspot.com",
    messagingSenderId: "1090069978845",
    appId: "1:1090069978845:web:f69c9f5aceac27eec20534",
    measurementId: "G-NXKVTQR1ZQ"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);

const auth = getAuth(app);

const googleProvider = new GoogleAuthProvider();

export { auth, signInAnonymously, googleProvider, signInWithPopup };