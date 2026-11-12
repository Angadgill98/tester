import React from 'react';
import { useGlobalStore } from '../../../global';
import HttpAction from './HttpAction/HttpAction';
import Sequence from './Sequence/Sequence';
import WsAction from './WsAction/wsAction';

const Dashboard = () => {
    const currentTab = useGlobalStore(state => state.CurrentTab);

    if (!currentTab) {
        return <div>No tab selected</div>;
    }

    return (
        <div>
            {currentTab.TabType === "http" && <HttpAction tab={currentTab} />}
            {currentTab.TabType === "sequence" && <Sequence tab={currentTab} />}
            {currentTab.TabType === "ws" && <WsAction tab={currentTab} />}
        </div>
    );
};

export default Dashboard;