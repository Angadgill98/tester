import React, { useContext, useEffect, useState } from 'react'
import { TesterSidebarClasses } from './Tester'
import { GetApplication,CreateHttpAction,CreateWSAction,GetApplications,CreateApplication,CreateSequence } from "../../../wailsjs/go/main/App";
import "./sidebar.css"
import { useGlobalStore } from '../../global';
import { Action_f, Application_f, ChildApplicationInfo, Environment_f, Sequence_f } from '../../models/application';


let ListTemplateClasses = {
    container: "sidebar-list-container",
    item: "sidebar-list-item",
    icon: "sidebarlisticon"
};
const Sidebar = () => {

    const setApplications = useGlobalStore(state => state.setApplications);
    const setActions = useGlobalStore(state => state.setActions);
    const setSequences = useGlobalStore(state => state.setSequences);

    const GlobalApplication = useGlobalStore(state => state.Applications.get(0) ?? null);
    
    const setCurrentApplicationId=useGlobalStore(state => state.setCurrentApplicationId);

    const currentApplicationId = useGlobalStore(state => state.CurrentApplicationId);

    const getCurrentApplication = useGlobalStore(state => state.getCurrentApplication);
    
    useEffect(() => {
        const fetchApplication = async () => {
            try {
                const visited = new Set<number>();
                let pendingIds = [0];

                const allApplications = new Map<number, Application_f>();
                const allActions = new Map<number, Action_f>();
                const allSequences = new Map<number, Sequence_f>();

                while (pendingIds.length > 0) {
                    pendingIds = [...new Set(pendingIds)].filter(id => !visited.has(id));

                    if (pendingIds.length === 0) break;

                    pendingIds.forEach(id => visited.add(id));

                    const res = await GetApplications(pendingIds);

                    if (!res.Success) {
                        console.error(res.Error);
                        return;
                    }

                    const applications = res.Data as {
                        Application: Application_f;
                        Actions: Action_f[];
                        Sequences: Sequence_f[];
                        Children: ChildApplicationInfo[];
                    }[];

                    const nextIds: number[] = [];

                    for (const item of applications) {
                        const application = item.Application;

                        item.Actions.forEach(action => allActions.set(action.Id, action));
                        item.Sequences.forEach(sequence => allSequences.set(sequence.Id, sequence));

                        item.Children.forEach(child => {
                            if (!visited.has(child.Id)) {
                                nextIds.push(child.Id);
                            }
                        });

                        allApplications.set(application.AppId, {
                            ...application,
                            Actions: item.Actions.map(action => action.Id),
                            Sequences: item.Sequences.map(sequence => sequence.Id),
                            Children: item.Children.map(child => child.Id)
                        });
                    }

                    pendingIds = nextIds;
                }

                // Publish everything to Zustand after fetching is complete.
                const store = useGlobalStore.getState();

                allActions.forEach((action, id) => store.setActions(id, action));
                allSequences.forEach((sequence, id) => store.setSequences(id, sequence));
                allApplications.forEach((application, id) => store.setApplications(id, application));

                store.setCurrentApplicationId(allApplications.get(0)?.AppId ?? null);

                console.log("Applications:", useGlobalStore.getState().Applications);

            } catch (error) {
                console.error("Failed to initialize applications:", error);
            }
        };

        fetchApplication();
    }, []);


    let [showActions, setActionShow] = useState<boolean>(false);
    let [showApplications, setApplicationShow] = useState<boolean>(false);
    let [showSequences, setSequenceShow] = useState<boolean>(false);

    const addActionToStore = (applicationId: number, actionId: number, actionName: string, actionType: string) => {
        const store = useGlobalStore.getState();
        const application = store.Applications.get(applicationId);

        if (!application) return;

        setActions(actionId, {
            Id: actionId,
            Name: actionName,
            Type: actionType
        });

        const updatedApplication = {
            ...application,
            Actions: [...application.Actions, actionId]
        };

        setApplications(applicationId, updatedApplication);

       
    };

    const handleCreateAction = () => {
        setActionShow(prev => !prev);
    };


    const handleCreateApplication = async () => {
        if (currentApplicationId === null) return;

        try {
            const response = await CreateApplication("New Application", currentApplicationId);

            if (!response.Success) {
                console.error(response.Error);
                return;
            }

            const childApp: Application_f = response.Data;

            const store = useGlobalStore.getState();
            const parentApp = store.Applications.get(currentApplicationId);

            if (!parentApp) return;

            setApplications(childApp.AppId, childApp);

            setApplications(currentApplicationId, {
                ...parentApp,
                Children: [...parentApp.Children, childApp.AppId]
            });

        } catch (error) {
            console.error("Failed to create application:", error);
        }
    };
    
    const handleCreateSequence = async () => {
        const store = useGlobalStore.getState();
        const currentApplication = store.getCurrentApplication();

        if (!currentApplication) return;

        const res = await CreateSequence("New Sequence", currentApplication.AppId);

        if (!res.Success) {
            console.error(res.Error);
            return;
        }

        const sequence = res.Data as Sequence_f;

        // Store the sequence object.
        store.setSequences(sequence.Id, sequence);

        // Add the sequence ID to its parent application.
        store.setApplications(currentApplication.AppId, {
            ...currentApplication,
            Sequences: [...currentApplication.Sequences, sequence.Id]
        });

        setSequenceShow(prev => !prev);
    };

    const handleCreateHttp = async () => {
        const currentApplication = useGlobalStore.getState().getCurrentApplication();

        if (!currentApplication) return;

        const applicationId = currentApplication.AppId;

        const res = await CreateHttpAction("New_Action", "http", applicationId);

        if (!res.Success) {
            console.error(res.Error);
            return;
        }

        addActionToStore(applicationId, res.Data, "New_Action", "http");

        console.log("Created action ID:", res.Data);
        setActionShow(false);
    };

    const handleCreateWs = async () => {
        const currentApplication = useGlobalStore.getState().getCurrentApplication();

        if (!currentApplication) return;

        const applicationId = currentApplication.AppId;

        const res = await CreateWSAction("New_Action", "ws", applicationId);

        if (!res.Success) {
            console.error(res.Error);
            return;
        }

        addActionToStore(applicationId, res.Data, "New_Action", "ws");

        console.log("Created WebSocket action ID:", res.Data);
        setActionShow(false);
    };

  return (
    <div className={TesterSidebarClasses.container}>

        <div className={TesterSidebarClasses.buttonsContainer}>

            <div className={TesterSidebarClasses.applicationButton}>
                <button onClick={handleCreateApplication}>Create Application</button>

                {/* {showApplications && (
                    <>
                        <div className={TesterSidebarClasses.buttonModaloverlay} onClick={() => setApplicationShow(false)} />
                        <div className="applcation-modal">
                            <button>Child Application</button>
                        </div>
                    </>
                )} */}
            </div>

            <div className="action-button">
                <button onClick={handleCreateAction}>Create Action</button>

                {showActions && (
                    <>
                        <div className={TesterSidebarClasses.buttonModaloverlay} onClick={() => setActionShow(false)} />
                        <div className="action-modal">
                            <button onClick={handleCreateHttp}>HTTP Action</button>
                            <button onClick={handleCreateWs}>WebSocket Action</button>
                        </div>
                    </>
                )}
            </div>

            <div className={TesterSidebarClasses.sequenceButton}>
                <button onClick={handleCreateSequence}>Create Sequence</button>

                {/* {showSequences && (
                    <>
                        <div className={TesterSidebarClasses.buttonModaloverlay} onClick={() => setSequenceShow(false)} />
                        <div className="sequence-modal">
                            <button>Create Sequence</button>
                        </div>
                    </>
                )} */}
            </div>

        </div>

        <div onClick={() => setCurrentApplicationId(GlobalApplication?.AppId ?? null)}>
            GLOBAL
        </div>

                {GlobalApplication ? (<SidebarTemplate application={GlobalApplication} classes={ListTemplateClasses} />) : (<></>)}

                
        
    </div>
  )
}

export default Sidebar




interface SidebarTemplateProps {
    application: Application_f;
    classes:any
}

const SidebarTemplate = React.memo(({ application, classes }: SidebarTemplateProps) => {

    
    const setCurrentApplicationId = useGlobalStore(state => state.setCurrentApplicationId);

    const setcurretnappid = () => {
        // console.log("new current id is ",application.AppId)
        setCurrentApplicationId(application.AppId);
    };

    const store = useGlobalStore.getState();

    return (
        <div className={classes.container} >
            {application.Children.map(childId => {
                const child = store.Applications.get(childId);
              

                if (!child) return null;

                return (
                    <SidebarApplicationTemplate key={child.AppId} application={child} />
                );
            })}

            {application.Actions.map(actionId => {
                const action = store.Actions.get(actionId);

                if (!action) return null;

                return <SidebarActionTemplate key={action.Id} action={action} updatecurrentappid={setcurretnappid} />;
            })}

            {application.Sequences.map(sequenceId => {
                const sequence = store.Sequences.get(sequenceId);

                if (!sequence) return null;

                return <SidebarSequenceTemplate key={sequence.Id} sequence={sequence} updatecurrentappid={setcurretnappid} />;
            })}
        </div>
    );
});

interface SidebarActionTemplateProps {
    action: Action_f;
    updatecurrentappid: () => void;
}

interface SidebarSequenceTemplateProps {
    sequence: Sequence_f;
    updatecurrentappid: () => void;
}

interface SidebarApplicationTemplateProps {
    application: Application_f;
}

const SidebarActionTemplate = React.memo(({ action, updatecurrentappid }: SidebarActionTemplateProps) => {
    const getCurrentApplicationId = useGlobalStore(state => state.getCurrentApplicationId);
    const addToTabs = useGlobalStore(state => state.addToTabs);

    const handleClick = () => {
        updatecurrentappid();

        const appId = getCurrentApplicationId();
        if (appId === null) return;

        addToTabs({
            TypeId: action.Id,
            AppId: appId,
            TabType: action.Type,
            Name:action.Name
        });
    };

    return (
        <div className={ListTemplateClasses.item} onClick={handleClick}>
            <div className={ListTemplateClasses.icon}>
                {action.Type.toUpperCase()}
            </div>
            <p>{action.Name}</p>
        </div>
    );
});

const SidebarSequenceTemplate = React.memo(({ sequence, updatecurrentappid }: SidebarSequenceTemplateProps) => {
    const getCurrentApplicationId = useGlobalStore(state => state.getCurrentApplicationId);
    const addToTabs = useGlobalStore(state => state.addToTabs);

    const handleClick = () => {
        updatecurrentappid();

        const appId = getCurrentApplicationId();
        if (appId === null) return;

        addToTabs({
            TypeId: sequence.Id,
            AppId: appId,
            TabType: "sequence",
            Name:sequence.Name
        });
    };
    return (
        <div className={ListTemplateClasses.item} onClick={handleClick}>
            <div className={ListTemplateClasses.icon}>S</div>
            <div>
                {/* <p>ID: {sequence.Id}</p> */}
                <p>Name: {sequence.Name}</p>
            </div>
        </div>
    );
});

const SidebarApplicationTemplate = React.memo(({ application }: SidebarApplicationTemplateProps) => {
    const [showDetails, setShowDetails] = useState<boolean>(false);

    const applicationData = useGlobalStore(state => state.Applications.get(application.AppId));
    const setCurrentApplicationId = useGlobalStore(state => state.setCurrentApplicationId);

    const handleClick = () => {
        if (!applicationData) return;

        setCurrentApplicationId(applicationData.AppId);
    };

    if (!applicationData) return null;

    return (
        <div>
            <div className={ListTemplateClasses.item}>
                <p onClick={handleClick}>Name: {applicationData.AppName}</p>

                <button type="button" onClick={() => {setShowDetails(prev => !prev);setCurrentApplicationId(applicationData.AppId);}}>
                    {showDetails ? "▲" : "▼"}
                </button>
            </div>

            {showDetails && (
                <SidebarTemplate application={applicationData} classes={ListTemplateClasses} />
            )}
        </div>
    );
});