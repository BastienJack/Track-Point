import {useState} from 'react'
import './login.css'
import { useNavigate } from 'react-router-dom'
import {serverUrl} from "../main.tsx";

function recordUserInfo(username) {
    localStorage.setItem('username', username)
}

export default function LoginPage() {
    const [username, setUsername] = useState("")
    function handleUsernameInput(event) {
        setUsername(event.target.value)
    }

    const [password, setPassword] = useState("")
    function handlePasswordInput(event) {
        setPassword(event.target.value)
    }

    const navigate = useNavigate()

    function login() {
        let loginPath = "/commerce/login"
        let url = serverUrl + loginPath

        let req = {"username": username, "password": password}
        
        fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(req)
        })
        .then(response => response.json())
        .then(data => {
            if (data['status_code'] != -1)
            {
                alert("Login success!")
                recordUserInfo(username)
            }
            navigate('/item-page')
        })
        .catch(error => {
            alert(error)
        })
    }

    function register() {
        let registerPath = "/commerce/register"
        let url = serverUrl + registerPath

        let req = {"username": username, "password": password, "confirm_password": password}

        fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(req)
        })
        .then(response => response.json())
        .then(data => {
            if (data['status_code'] != -1)
            {
                alert("Register success!")
            }
        })
        .catch(error => {
            alert(error)
        })
    }

    return (
        <div>
            <div className='background'/>

            <div className='hint'>
                Hello
            </div>
    
            <div className='input'>
                <div>
                    <input 
                    type='text'
                    className='username_input' 
                    placeholder='username'
                    value={username}
                    onChange={handleUsernameInput}
                    />
                </div>
                <div>
                    <input 
                    type='text'
                    className='password_input'
                    placeholder='password'
                    value={password}
                    onChange={handlePasswordInput}
                    />
                </div>
            </div>
    
            <div className='button'>
                <button className="login_button" onClick={login}>Login</button>
                <button className="register_button" onClick={register}>Register</button>
            </div>
        </div>
    )
}

