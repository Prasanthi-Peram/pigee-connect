import './Login.css';
import {useState} from 'react';
import axios from 'axios';


function Login(){
    const [email,setEmail]= useState('');

    const [password,setPassword]=useState('');

    const handleSubmit = async (e:React.FormEvent)=>{
        e.preventDefault();
        try{
            const res= await axios.post('http://localhost:8080/v1/authentication/token',{
                email,
                password
            });

            const token=res.data as string
            localStorage.setItem("authToken",token);
            console.log("Token daved:",token);
            console.log(res.data)
            alert('Login successful')
        }catch(err){
            console.log(err)
            alert('Login failed')
        }
    }

    return(
        <div className='login-container'>
        <form className='login-form' onSubmit={handleSubmit}>
            <h2>Login</h2>
            <input type="email" placeholder="Email" value={email} onChange={(e)=>setEmail(e.target.value)} required/>

            <input type="password" placeholder="Password"
            value={password}
            onChange={(e)=>setPassword(e.target.value)} required/>

            <button type="submit">Login</button>
            </form>
    </div>
    );
  
}

export default Login;