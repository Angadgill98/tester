import React, { useContext, useEffect, useState } from 'react'
import { TesterSidebarClasses } from './Tester'
import { GetApplication,CreateHttpAction,CreateWSAction,GetApplications } from "../../../wailsjs/go/main/App";
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
    const setChildApplications = useGlobalStore(state => state.setChildApplications);

    const GlobalApplication = useGlobalStore(state => state.GlobalApplication);

    const setGlobalApplication = useGlobalStore(state => state.setGlobalApplication);
    
    const setCurrentApplicationId=useGlobalStore(state => state.setCurrentApplicationId);
    
    useEffect(() => {
        const fetchApplication = async () => {
            try {
                const visited = new Set<number>();
                let pendingIds = [0];

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

                        // Store actions separately.
                        for (const action of item.Actions) {
                            setActions(action.Id, action);
                        }

                        // Store sequences separately.
                        for (const sequence of item.Sequences) {
                            setSequences(sequence.Id, sequence);
                        }

                        // Store child application information separately.
                        for (const child of item.Children) {
                            setChildApplications(child.Id, child);

                            if (!visited.has(child.Id)) {
                                nextIds.push(child.Id);
                            }
                        }

                        // Store only IDs inside the application.
                        const applicationData: Application_f = {
                            ...application,
                            Actions: item.Actions.map(action => action.Id),
                            Sequences: item.Sequences.map(sequence => sequence.Id),
                            Children: item.Children.map(child => child.Id)
                        };

                        setApplications(applicationData.AppId, applicationData);
                    }

                    pendingIds = nextIds;
                }

                const store = useGlobalStore.getState();
                const rootApplication = store.Applications.get(0) ?? null;

                store.setGlobalApplication(rootApplication);
                store.setCurrentApplicationId(rootApplication?.AppId ?? null);

                console.log("Applications:", store.Applications);
                console.log("Actions:", store.Actions);
                console.log("Sequences:", store.Sequences);
                console.log("Child Applications:", store.ChildApplications);

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

        if (store.GlobalApplication?.AppId === applicationId) {
            setGlobalApplication(updatedApplication);
        }
    };
  const handleCreateAction = () => {
        setActionShow(prev => !prev);
    };

  const handleCreateApplication = () => {
      setApplicationShow(prev => !prev);
  };

  const handleCreateSequence = () => {
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

                {showApplications && (
                    <>
                        <div className={TesterSidebarClasses.buttonModaloverlay} onClick={() => setApplicationShow(false)} />
                        <div className="applcation-modal">
                            <button>Child Application</button>
                        </div>
                    </>
                )}
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

                {showSequences && (
                    <>
                        <div className={TesterSidebarClasses.buttonModaloverlay} onClick={() => setSequenceShow(false)} />
                        <div className="sequence-modal">
                            <button>Create Sequence</button>
                        </div>
                    </>
                )}
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
        setCurrentApplicationId(application.AppId);
    };

    const store = useGlobalStore.getState();

    return (
        <div className={classes.container} >
            {application.Children.map(childId => {
                const child = store.ChildApplications.get(childId);

                if (!child) return null;

                return (
                    <SidebarApplicationTemplate key={child.Id} application={child} />
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
    application: ChildApplicationInfo;
}

const SidebarActionTemplate = React.memo(({ action, updatecurrentappid }: SidebarActionTemplateProps) => {
    return (
        <div className={ListTemplateClasses.item} onClick={updatecurrentappid}>
            <div className={ListTemplateClasses.icon}>
                {action.Type.toUpperCase()}
            </div>
            <p>{action.Name}</p>
        </div>
    );
});

const SidebarSequenceTemplate = React.memo(({ sequence, updatecurrentappid }: SidebarSequenceTemplateProps) => {
    return (
        <div className={ListTemplateClasses.item} onClick={updatecurrentappid}>
            <div className={ListTemplateClasses.icon}>S</div>
            <div>
                <p>ID: {sequence.Id}</p>
                <p>Name: {sequence.Name}</p>
            </div>
        </div>
    );
});

const SidebarApplicationTemplate = React.memo(({ application }: SidebarApplicationTemplateProps) => {
    const [showDetails, setShowDetails] = useState<boolean>(false);

    const applicationData = useGlobalStore(state => state.Applications.get(application.Id));
    const setCurrentApplicationId = useGlobalStore(state => state.setCurrentApplicationId);

    const handleClick = () => {
        if (!applicationData) return;

        setCurrentApplicationId(applicationData.AppId);
        setShowDetails(prev => !prev);
    };

    if (!applicationData) return null;

    return (
        <div onClick={handleClick}>
            <div className={ListTemplateClasses.item} >
                <p>Name: {application.Name}</p>
                <button type="button">
                    {showDetails ? "▲" : "▼"}
                </button>
            </div>

            {showDetails && (
                <SidebarTemplate
                    application={applicationData}
                    classes={ListTemplateClasses}
                />
            )}
        </div>
    );
});