import React, { useEffect, useState } from 'react';
import { Response, Tab, WSAction as WSActionType } from '../../../../models/application';
import { invoke } from '@tauri-apps/api/core';

interface WsActionProps {
    tab: Tab;
}

const WsAction = ({ tab }: WsActionProps) => {
    const [wsAction, setWsAction] = useState<WSActionType | null>(null);

    useEffect(() => {
        let cancelled = false;

        setWsAction(null);

        const fetchWSAction = async () => {
            try {
              const response = await invoke<Response<WSActionType>>("get_ws_action", {
                  actionId: tab.TypeId
              });
              
              if (cancelled) return;
              
              if (!response.success) {
                  console.error(response.error);
                  return;
              }

                setWsAction(response.data);
            } catch (error) {
                if (!cancelled) {
                    console.error("Failed to fetch WebSocket action:", error);
                }
            }
        };

        fetchWSAction();

        return () => {
            cancelled = true;
        };
    }, [tab.TypeId, tab.AppId, tab.TabType]);

    if (!wsAction) return <div>Loading WebSocket action...</div>;

    return (
        <div>
            <h2>{wsAction.ActionName}</h2>
            <p>WebSocket URL: {wsAction.Opts.URL}</p>
            <p>Action ID: {wsAction.ActionID}</p>
            <p>App ID: {wsAction.AppID}</p>
        </div>
    ); 
};

export default WsAction;