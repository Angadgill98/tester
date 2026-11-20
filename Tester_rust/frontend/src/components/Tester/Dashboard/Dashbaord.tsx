import React from 'react';
import { useGlobalStore } from '../../../global';
import HttpAction from './HttpAction/HttpAction';
import Sequence from './Sequence/Sequence';
import WsAction from './WsAction/wsAction';
import "./dashboard.css"
import { DashboardClasses } from '../Tester';


export let HttpClasses={
  container: "http-container",
  action: "http-action",
  actionInfo: "http-action-info",
  main: "http-main",
  main_left: "http-main-left",
  main_right: "http-main-right",



  http_opts: "http-opts",
  http_opts_tab: "http-opts-tab",
  http_opts_field:"http-opts-field",
  http_rows_first: "http_rows_first",
  http_rows_second: "http_rows_second",
  http_rows_third:"http_rows_third",


  http_headers_rows: "http-headers-row",
  http_queries_rows:"http-queries-row",
  
  
  
}
const Dashboard = () => {
    const currentTab = useGlobalStore(state => state.CurrentTab);

    if (!currentTab) {
        return <div>No tab selected</div>;
    }

    return (
        <div className={DashboardClasses.container}>
            {currentTab.TabType === "http" && <HttpAction tab={currentTab} />}
            {currentTab.TabType === "sequence" && <Sequence tab={currentTab} />}
            {currentTab.TabType === "ws" && <WsAction tab={currentTab} />}
        </div>
    );
};

export default Dashboard;