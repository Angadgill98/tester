import React, { useEffect, useState } from 'react';
import { GetWSAction } from "../../../../../wailsjs/go/main/App";
import { Tab, WSAction as WSActionType } from '../../../../models/application';

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
                const response = await GetWSAction(tab.TypeId);

                if (cancelled) return;

                if (!response.Success) {
                    console.error(response.Error);
                    return;
                }

                setWsAction(response.Data);
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