import React from 'react'
import { TesterSidebarClasses } from './Tester'
import "./sidebar.css"
const Sidebar = () => {
  return (
    <div className={TesterSidebarClasses.container}>

        <div>
            <button>Create Action</button>
            <button>Create Sequence</button>
            <button>Create Application</button>
        </div>

        <div>GLOBAL</div>


        
    </div>
  )
}

export default Sidebar


function GetAppDetails(){
    
}