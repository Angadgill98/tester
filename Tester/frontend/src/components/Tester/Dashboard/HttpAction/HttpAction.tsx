import React, { useEffect, useState } from 'react';
import { GetHttpAction } from "../../../../../wailsjs/go/main/App";
import { Tab, HttpAction as HttpActionType } from '../../../../models/application';

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
    }, [tab]);

    if (!action) return <div>Loading HTTP action...</div>;

    return (
        <div>
            <h2>{action.ActionName}</h2>
            <p>WebSocket URL: {action.Opts.URL}</p>
            <p>Action ID: {action.ActionId}</p>
            <p>App ID: {action.AppId}</p>
        </div>
    );
};

export default HttpAction;