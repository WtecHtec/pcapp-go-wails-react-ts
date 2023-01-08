import {useState,useEffect} from 'react';
// import logo from './assets/images/logo-universal.png';
import { Button } from 'antd';
import './App.css';
import {LoginCode, GetQRpath} from "../wailsjs/go/main/App";
import { EventsOn, EventsOff } from '../wailsjs/runtime/runtime'
import { useNavigate  } from 'react-router-dom'

function App() {
    const navigate = useNavigate()
    const [loginCode, setLoginCode] = useState('build/appicon.png')
    const [status, setStatus] = useState(false);
 
    useEffect(() => {
      const login = async (unid: string) => {
        if (unid === '200') {
          navigate('/home')
          return;
        }
        const qrPath = await GetQRpath()
        setLoginCode(`${qrPath}${unid}.png`)
        setStatus(true)
      }
      EventsOn('login:code', login)
      return () => {
        EventsOff('login:code', 'login')
      }
    }, [])

    function handleLogin() {
        LoginCode()
    }
    return (
        <div className="app">
            <span className="c-333">长时间不响应,可尝试刷新登录二维码</span>
            <img src={loginCode} className="logo" alt="登录二维码"/>
            <div  className="m-t-24">
              <Button type="primary" onClick={handleLogin}> { status ? '刷新二维码' : '登录'}</Button>
            </div>
        </div>
    )
}

export default App
