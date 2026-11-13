import React, { useEffect, useState } from "react";
import { GetHttpAction } from "../../../../../wailsjs/go/main/App";
import {
  Tab,
  HttpAction as HttpActionType,
} from "../../../../models/application";
import { HttpClasses } from "../Dashbaord";
import "./Httpaction.css";

interface HttpActionProps {
  tab: Tab;
}

const HttpAction = ({ tab }: HttpActionProps) => {
  const [action, setAction] = useState<HttpActionType | null>(null);
  useEffect(() => {
    let cancelled = false;

    const fetchAction = async () => {
      try {
        const response = await GetHttpAction(tab.TypeId);

        if (cancelled) return;

        if (!response.Success) {
          console.error(response.Error);
          setAction(null);
          return;
        }
        console.log("Got the HttpAction: ", response.Data);
        setAction(response.Data);
      } catch (error) {
        if (cancelled) return;

        console.error(error);
        setAction(null);
      }
    };

    fetchAction();

    return () => {
      cancelled = true;
    };
  }, [tab.TypeId]);

  if (!action) return <div>Loading HTTP action...</div>;

  return (
    <div className={HttpClasses.container}>
      <div className={HttpClasses.action}>
        <div className={HttpClasses.actionInfo}>
          <h3>Action Name</h3>
          <span>{action.ActionName}</span>
        </div>

        <div className={HttpClasses.actionInfo}>
          <p>Application Name</p>
        </div>

        <div className={HttpClasses.actionInfo}>
          <p>Environment</p>
        </div>
      </div>


      <div className={HttpClasses.main}>
        <div className={HttpClasses.main_left}>
          <p>Http Options</p>
    
          <div>
            <div>req</div>
    
            <div>settings</div>
    
            <div></div>
          </div>

        </div>

        <div className={HttpClasses.main_right}>
            <p>Response</p>
        </div>
        
      </div>
      
     

      <div>
        <p>History</p>
      </div>
    </div>
  );
};

export default HttpAction;
