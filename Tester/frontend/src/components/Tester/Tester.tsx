import React from 'react'
import { TesterClasses } from '../../App'
import "./tester.css" 
import Sidebar from './Sidebar'
import Tabs from './Tabs/Tabs'
import Dashboard from './Dashboard/Dashbaord'

export let TesterSidebarClasses={
    container:"tester-sidebar-container",
    buttonsContainer:"button-container",
    actionButton:"action-button",
    sequenceButton:"sequence-button",
    applicationButton:"application-button",
    buttonModaloverlay:"modal-overlay",
    
}

export let TabsClasses={
  container:"tab-container"
}

export let DashboardClasses = {
  container: "dashboard-container",
  
}

const Tester = () => {
  return (
    <div className={TesterClasses.container}>
        <Sidebar/>
        <div className={TesterClasses.right}>
          <Tabs />
          <Dashboard/>
        </div>
        
    </div>
  )
}

export default Tester