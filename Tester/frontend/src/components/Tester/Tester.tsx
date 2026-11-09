import React from 'react'
import { TesterClasses } from '../../App'
import "./tester.css" 
import Sidebar from './Sidebar'

export let TesterSidebarClasses={
    container:"tester-sidebar-container"
}

const Tester = () => {
  return (
    <div className={TesterClasses.container}>
        <Sidebar/>
        asdasd
    </div>
  )
}

export default Tester