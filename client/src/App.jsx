import Navbar from "./components/navbar/Navbar"
import './App.css'
import { Routes, Route } from "react-router-dom"
import Home from "./sites/home/Home"
import Register from "./sites/register/Register"
import Login from "./sites/login/Login"

function App() {

  return (
    <section className="App">
      <Navbar></Navbar>
      <Routes>
        <Route path='/' element={<Home/>}/>
        <Route path='/register' element={<Register/>}/>
        <Route path='/login' element={<Login/>}/>
      </Routes>
    </section>
  )
}

export default App
