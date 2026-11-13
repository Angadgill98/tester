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

  let [showreq, setshow] = useState<boolean>(true)
  let [showbody, setshowbody] = useState<boolean>(false)
  
  
  let handleReqclick = () => {
    setshow(true)
    setshowbody(false)
    
  }

  let handleBodyclick = () => {
    setshowbody(true)
    setshow(false)
  }

  
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

        <button>Settings</button>
        <button>Toggle Temp Sequnce</button>
        <button>Reset temp Sequnce</button>
        <button>Execute</button>
      </div>


      <div className={HttpClasses.main}>
        <div className={HttpClasses.main_left}>
          <p>Http Options</p>
    
          <div className={HttpClasses.http_opts}>
            <div className={HttpClasses.http_opts_tab}>
              <button onClick={handleReqclick}>req</button>
              <button onClick={handleBodyclick}>body</button>
            </div>

            {showreq && (
                <div>
                  <div>
                      <label>Method</label>
                      <select defaultValue={action.Opts.Method}>
                          <option value="GET">GET</option>
                          <option value="POST">POST</option>
                          <option value="PUT">PUT</option>
                          <option value="PATCH">PATCH</option>
                          <option value="DELETE">DELETE</option>
                          <option value="HEAD">HEAD</option>
                      </select>
                  </div>
            
                    <div>
                        <label>URL</label>
                        <input type="text" defaultValue={action.Opts.URL} />
                    </div>
            
                    <div>
                        <label>Headers</label>
                    
                        {Object.entries(action.Opts.Headers).map(([key, value]) => (
                            <div key={key}>
                                <input type="text" defaultValue={key} placeholder="Key" />
                                <input type="text" defaultValue={String(value)} placeholder="Value" />
                            </div>
                        ))}
                    </div>
                    
                    <div>
                        <label>Query</label>
                    
                        {Object.entries(action.Opts.Query).map(([key, value]) => (
                            <div key={key} >
                                <input type="text" defaultValue={key} placeholder="Key" />
                                <input type="text" defaultValue={String(value)} placeholder="Value" />
                            </div>
                        ))}
                    </div>
            
                    <div>
                        <label>Cookies</label>
                    
                        {Object.entries(action.Opts.Cookies).map(([key, value]) => (
                            <div key={key} style={{ display: "flex", gap: "10px" }}>
                                <input type="text" defaultValue={key} placeholder="Cookie Name" />
                                <input type="text" defaultValue={String(value)} placeholder="Cookie Value" />
                            </div>
                        ))}
                    </div>
                </div>
            )}
            
            {showbody && (
                <div>
                    <label>Body</label>
                    <textarea defaultValue={JSON.stringify(action.Opts.Body, null, 2)} />
                </div>
            )}
    
    
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
