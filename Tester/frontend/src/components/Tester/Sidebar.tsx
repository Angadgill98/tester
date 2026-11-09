import React, { useEffect } from 'react'
import { TesterSidebarClasses } from './Tester'
import { GetApplication } from "../../../wailsjs/go/main/App";
import "./sidebar.css"
const Sidebar = () => {
  useEffect(() => {
      const fetchApplication = async () => {
          try {
              const application = await GetApplication(0);
              console.log(application);
          } catch (error) {
              console.error("Failed to get application:", error);
          }
      };

      fetchApplication();
  }, []);

  
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


async function GetAppDetails(){
  let app=await GetApplication(0)
  console.log(app)
}