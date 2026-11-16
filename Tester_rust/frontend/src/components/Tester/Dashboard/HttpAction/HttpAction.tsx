import React, { useEffect, useMemo, useState } from "react";
import {
  Tab,
  HttpAction as HttpActionType,
  Response,
} from "../../../../models/application";
import { HttpClasses } from "../Dashbaord";
import "./Httpaction.css";
import { invoke } from "@tauri-apps/api/core";


interface KeyValuePair {
  id: number;
  key: string;
  value: string;
}

interface HttpActionProps {
  tab: Tab;
}

const HttpAction = ({ tab }: HttpActionProps) => {
  const [action, setAction] = useState<HttpActionType | null>(null);
  const [headers, setHeaders] = useState<KeyValuePair[]>([]);
  const [query, setQuery] = useState<KeyValuePair[]>([]);
  const [cookie,setCookie]=useState<KeyValuePair[]>([]);

  const [preview_req, setPreview] = useState<string>("")
  const [curl_req,setCurl]=useState<string>("")
  
  useEffect(() => {
    let cancelled = false;
  
    const fetchAction = async () => {
      try {
        const response = await invoke<Response<HttpActionType>>("get_http_action", {
            actionId: tab.TypeId
        });
        
        if (cancelled) return;
        
        if (!response.success) {
            console.error(response.error);
            setAction(null);
            return;
        }
  
        console.log("Got the HttpAction: ", response.data);
        setAction(response.data);

        setHeaders(
            Object.entries(response.data.Opts.Headers).flatMap(([key, values]) =>
                values.map(value => ({
                    id: Date.now() + Math.random(),
                    key,
                    value
                }))
            )
        );
    
        setQuery(
            Object.entries(response.data.Opts.Query).flatMap(([key, values]) =>
                values.map(value => ({
                    id: Date.now() + Math.random(),
                    key,
                    value
                }))
            )
        );

        setCookie(
            Object.entries(response.data.Opts.Cookies).flatMap(([key, values]) =>
                values.map(value => ({
                    id: Date.now() + Math.random(),
                    key,
                    value
                }))
            )
        );
  
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

  useEffect(() => {
      if (!action) {
          return;
      }
    let preview = getPreviewUrl();
    let curl = getCurlRequest();

    setCurl(curl)
    setPreview(preview)
    
    
  }, [action]);

  let handleUrlChanege = (newurl:string) => {
    const updatedAction = {
        ...action!,
        Opts: {
            ...action!.Opts,
            Url: newurl
        }
    };

    setAction(updatedAction);
    debouncedSaveAction(updatedAction);
  }

  const handleMethodChange = (newMethod: string) => {
      const updatedAction = {
          ...action!,
          Opts: {
              ...action!.Opts,
              Method: newMethod
          }
      };
  
      setAction(updatedAction);
      debouncedSaveAction(updatedAction);
  };


  let [showreq, setshow] = useState<boolean>(true);
  let [showbody, setshowbody] = useState<boolean>(false);
  let [showcookie, setshowcookie] = useState<boolean>(false);
  
  let handleReqclick = () => {
    setshow(true);
    setshowbody(false);
    setshowcookie(false);
  };
  
  let handleBodyclick = () => {
    setshowbody(true);
    setshow(false);
    setshowcookie(false);
  };
  
  let handleCookieclick = () => {
    setshowcookie(true);
    setshow(false);
    setshowbody(false);
  };

  const handleAddHeader = async () => {
      const new_key_value_pair: KeyValuePair = {
          id: Date.now(),
          key: "key",
          value: "value"
      };
  
      const newHeaders = convertToMap([...headers, new_key_value_pair]);
  
      const updatedAction = {
          ...action!,
          Opts: {
              ...action!.Opts,
              Headers: newHeaders
          }
      };
  
      const response = await invoke<Response<null>>("save_http_action", {
          action: updatedAction
      });

    console.log(response)
  
      if (!response.success) {
          console.error(response.error);
          return;
      }
      console.log(updatedAction)
      
      setAction(updatedAction);
      setHeaders(prev => [...prev, new_key_value_pair]);
  };
  
  const handleAddQuery = async () => {
      const new_key_value_pair: KeyValuePair = {
          id: Date.now(),
          key: "key",
          value: "value"
      };
  
      const newQuery = convertToMap([...query, new_key_value_pair]);
  
      const updatedAction = {
          ...action!,
          Opts: {
              ...action!.Opts,
              Query: newQuery
          }
      };
  
      const response = await invoke<Response<null>>("save_http_action", {
          action: updatedAction
      });
  
      if (!response.success) {
          console.error(response.error);
          return;
      }
  
      setAction(updatedAction);
      setQuery(prev => [...prev, new_key_value_pair]);
  };

  const handleAddCookie = async () => {
      const new_key_value_pair: KeyValuePair = {
          id: Date.now(),
          key: "key",
          value: "value"
      };
  
      const newCookies = convertToMap([...cookie, new_key_value_pair]);
  
      const updatedAction = {
          ...action!,
          Opts: {
              ...action!.Opts,
              Cookies: newCookies
          }
      };
  
      const response = await invoke<Response<null>>("save_http_action", {
          action: updatedAction
      });
  
      if (!response.success) {
          console.error(response.error);
          return;
      }
      console.log("asds  ",updatedAction)
      
      setAction(updatedAction);
      setCookie(prev => [...prev, new_key_value_pair]);
  };

  const debounce = (callback: Function, delay: number) => {
      let timer: ReturnType<typeof setTimeout>;
  
      return (...args: any[]) => {
          clearTimeout(timer);
          timer = setTimeout(() => callback(...args), delay);
      };
  };

  const debouncedSaveAction = useMemo(
      () => debounce(async (updatedAction: HttpActionType) => {
          const response = await invoke<Response<null>>("save_http_action", {
              action: updatedAction
          });
  
          if (!response.success) {
              console.error(response.error);
          }
      }, 500),
      []
  );
  
  const handleHeaderChange = (id: number, field: "key" | "value", value: string) => {
      const newHeaders = headers.map(item => item.id === id ? { ...item, [field]: value } : item);
  
      const updatedAction = {
          ...action!,
          Opts: {
              ...action!.Opts,
              Headers: convertToMap(newHeaders)
          }
      };
  
      setHeaders(newHeaders);
      setAction(updatedAction);
  
      debouncedSaveAction(updatedAction);
  };

  const handleQueryChange = (id: number, field: "key" | "value", value: string) => {
      const newQuery = query.map(item => item.id === id ? { ...item, [field]: value } : item);
  
      const updatedAction = {
          ...action!,
          Opts: {
              ...action!.Opts,
              Query: convertToMap(newQuery)
          }
      };
  
      setQuery(newQuery);
      setAction(updatedAction);
  
      debouncedSaveAction(updatedAction);
  };

  const handleCookieChange = (id: number, field: "key" | "value", value: string) => {
      const newCookies = cookie.map(item => item.id === id ? { ...item, [field]: value } : item);
  
      const updatedAction = {
          ...action!,
          Opts: {
              ...action!.Opts,
              Cookies: convertToMap(newCookies)
          }
      };

      console.log("asds  ",updatedAction)
      setCookie(newCookies);
      setAction(updatedAction);
      debouncedSaveAction(updatedAction);
  };

  const convertToMap = (items: KeyValuePair[]): Record<string, string[]> => {
      const result: Record<string, string[]> = {};
  
      items.forEach(({ key, value }) => {
          if (!key.trim()) return;
  
          if (!result[key]) {
              result[key] = [];
          }
  
          result[key].push(value);
      });
  
      return result;
  };

  const handleRemoveItem = async (id: number, type: "header" | "query" | "cookie") => {
      const items = type === "header" ? headers : type === "query" ? query : cookie;
      const updatedItems = items.filter(item => item.id !== id);
  
      const updatedAction = {
          ...action!,
          Opts: {
              ...action!.Opts,
              ...(type === "header" && { Headers: convertToMap(updatedItems) }),
              ...(type === "query" && { Query: convertToMap(updatedItems) }),
              ...(type === "cookie" && { Cookies: convertToMap(updatedItems) })
          }
      };
  
      const response = await invoke<Response<null>>("save_http_action", {
          action: updatedAction
      });
  
      if (!response.success) {
          console.error(response.error);
          return;
      }
  
      setAction(updatedAction);
  
      if (type === "header") setHeaders(updatedItems);
      if (type === "query") setQuery(updatedItems);
      if (type === "cookie") setCookie(updatedItems);
  };

  const getPreviewUrl = () => {
      if (!action) return "";
  
      const url = action.Opts.Url;
      const queryParams = new URLSearchParams();
  
      Object.entries(action.Opts.Query).forEach(([key, values]) => {
          values.forEach(value => queryParams.append(key, value));
      });
  
      const queryString = queryParams.toString();
  
      return queryString ? `${url}?${queryString}` : url;
  };
  
  const getCurlRequest = () => {
      if (!action) return "";
  
      const method = action.Opts.Method;
      const url = getPreviewUrl();
  
      const headers = Object.entries(action.Opts.Headers)
          .flatMap(([key, values]) =>
              values.map(value => `-H '${key}: ${value}'`)
          );
  
      const cookies = Object.entries(action.Opts.Cookies)
          .flatMap(([key, values]) =>
              values.map(value => `${key}=${value}`)
          )
          .join("; ");
  
      const curlParts = [
          `curl -X ${method}`,
          ...headers,
          ...(cookies ? [`-H 'Cookie: ${cookies}'`] : []),
          `'${url}'`
      ];
  
      return curlParts.join(" \\\n  ");
  };

  const handleCopyPreviewUrl = async () => {
      try {
          await navigator.clipboard.writeText(preview_req);
      } catch (error) {
          console.error("Failed to copy preview URL:", error);
      }
  };
  
  const handleCopyCurlRequest = async () => {
      try {
          await navigator.clipboard.writeText(curl_req);
      } catch (error) {
          console.error("Failed to copy cURL request:", error);
      }
  };
 
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
              <button onClick={handleCookieclick}>cookie</button>
            </div>

            {showreq && (
              <div style={{height:"95%"}}>
                <div className={HttpClasses.http_rows_first}>
                  <div className={HttpClasses.http_opts_field}>
                    <label style={{paddingLeft:"5px",boxSizing:"border-box"}}>Method</label>
                    <select value={action.Opts.Method} onChange={(e)=>handleMethodChange(e.target.value)}>
                      <option value="GET">GET</option>
                      <option value="POST">POST</option>
                      <option value="PUT">PUT</option>
                      <option value="PATCH">PATCH</option>
                      <option value="DELETE">DELETE</option>
                      <option value="HEAD">HEAD</option>
                    </select>
                  </div>

                  <div className={HttpClasses.http_opts_field} style={{flex:1}}>
                    <label style={{paddingLeft:"5px",boxSizing:"border-box"}}>URL</label>
                    <input type="text" onChange={(e)=>{handleUrlChanege(e.target.value)}} value={action.Opts.Url} />
                  </div>
                </div>

                <div className={HttpClasses.http_rows_second}>
                  <div className={HttpClasses.http_opts_field} style={{width:"50%"}}>
                    <div>
                      <label style={{ paddingLeft: "5px", boxSizing: "border-box" }}>Headers</label>
                      <button onClick={handleAddHeader}>Add</button>
                    </div >

                    <div className={HttpClasses.http_headers_rows}>
                      {headers.map(header => (
                          <div key={header.id}>
                              <input
                                  value={header.key}
                                  placeholder="Key"
                                  onChange={e => handleHeaderChange(header.id, "key", e.target.value)}
                              />
                      
                              <input
                                  value={header.value}
                                  placeholder="Value"
                                  onChange={e => handleHeaderChange(header.id, "value", e.target.value)}
                          />

                          <button onClick={() => handleRemoveItem(header.id, "header")}>Remove</button>
                          </div>
                      ))}
                   
                    </div>
                    
                  </div>
  
                  <div className={HttpClasses.http_opts_field} style={{width:"50%"}}>
                    <div>
                      <label style={{ paddingLeft: "5px", boxSizing: "border-box" }}>Query</label>
                      <button onClick={handleAddQuery}>Add</button>
                    </div>

                    <div className={HttpClasses.http_queries_rows}>
                      {query.map(item => (
                        <div key={item.id}>
                          <input
                            value={item.key}
                            placeholder="Key"
                            onChange={e => handleQueryChange(item.id, "key", e.target.value)}
                          />
                      
                          <input
                            value={item.value}
                            placeholder="Value"
                            onChange={e => handleQueryChange(item.id, "value", e.target.value)}
                          />

                          <button onClick={() => handleRemoveItem(item.id, "query")}>Remove</button>
                        </div>
                      ))}
                  
                    </div>
                   
                  </div>
  
                 
                </div>


                <div className={HttpClasses.http_rows_third}>
                    <div className={HttpClasses.http_opts_field}>
                        preview url
                        <span>{preview_req}</span>
                        <button onClick={handleCopyPreviewUrl}>Copy</button>
                    </div>
                
                    <div className={HttpClasses.http_opts_field}>
                        curl req
                        <span>{curl_req}</span>
                        <button onClick={handleCopyCurlRequest}>Copy</button>
                    </div>
                </div>
                
              </div>
            )}

            {showbody && (
              <div>
                <label>Body</label>
                <textarea
                  defaultValue={JSON.stringify(action.Opts.Body, null, 2)}
                />
              </div>
            )}

            {showcookie && (
              <div className={HttpClasses.http_opts_field}>
                  <label>Cookies</label>
              
                  {cookie.map(item => (
                      <div key={item.id}>
                          <input
                              value={item.key}
                              placeholder="Cookie Name"
                              onChange={e => handleCookieChange(item.id, "key", e.target.value)}
                          />
              
                          <input
                              value={item.value}
                              placeholder="Cookie Value"
                              onChange={e => handleCookieChange(item.id, "value", e.target.value)}
                      />
                      <button onClick={() => handleRemoveItem(item.id, "cookie")}>Remove</button>
                      </div>
                  ))}
              
                  <button onClick={handleAddCookie}>Add Cookie</button>
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
