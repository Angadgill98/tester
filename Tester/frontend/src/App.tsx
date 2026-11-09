import {useState} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet} from "../../../Tester/frontend/wailsjs/go/main/App";
import MainSideBar from './components/Sidebar/MainSideBar';
import Tester from './components/Tester/Tester';


export let SidebarClasses = {
    container:"sidebar-container",
    container_closed:"sidebar-container-closed",
    trigger:"sidebar-trigger",
    button:"sidebar-button",
    
}

export let TesterClasses={
    container:"tester-container"
}

function App() {
    let [isTester,setTester]=useState<boolean>(true)

    return (
        <div id="App">
            <MainSideBar/>
            {isTester ? <Tester/> :<></>}
        </div>
    )
}

export default App
