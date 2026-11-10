import React, { useContext, useEffect, useState } from 'react'
import { TesterSidebarClasses } from './Tester'
import { GetApplication,CreateHttpAction,CreateWSAction,GetApplications } from "../../../wailsjs/go/main/App";
import "./sidebar.css"
import { useGlobalStore } from '../../global';
import { Action_f, Application_f, ChildApplicationInfo, Environment_f, Sequence_f } from '../../models/application';
const Sidebar = () => {

    const { setApplications, setActions, setSequences, setChildApplications ,CurrentApplication} = useGlobalStore()

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
                store.setCurrentApplication(rootApplication);

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


  const handleCreateAction = () => {
        setActionShow(prev => !prev);
    };

  const handleCreateApplication = () => {
      setApplicationShow(prev => !prev);
  };

  const handleCreateSequence = () => {
      setSequenceShow(prev => !prev);
  };

  let handleCreateHttp = async () => {
    const res = await CreateHttpAction("New_Action", "http", CurrentApplication!.AppId);

    if (!res.Success) {
        console.error(res.Error);
        return;
    }

    console.log("Created action ID:", res.Data);
    setActionShow(prev => !prev);


  };

  let handleCreateWs=async ()=>{
    const res = await CreateWSAction("New_Action", "ws", CurrentApplication!.AppId);

    if (!res.Success) {
        console.error(res.Error);
        return;
    }

    console.log("Created WebSocket action ID:", res.Data);
    setActionShow(prev => !prev);
  }

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

        <div>GLOBAL</div>


        
    </div>
  )
}

export default Sidebar




interface SidebarTemplateProps {
    application: Application_f;
}

const SidebarTemplate = ({ application }: SidebarTemplateProps) => {
    const ChildApplications = useGlobalStore(state => state.ChildApplications);
    const Actions = useGlobalStore(state => state.Actions);
    const Sequences = useGlobalStore(state => state.Sequences);
    return (
        <div>
            {application.Children.map(childId => {
            const child = ChildApplications.get(childId);

            if (!child) return null;

            return (
                <SidebarApplicationTemplate key={child.Id} application={child} />
            );
        })}

        {application.Actions.map(actionId => {
            const action = Actions.get(actionId);

            if (!action) return null;

            return (
                <SidebarActionTemplate key={action.Id} action={action} />
            );
        })}

        {application.Sequences.map(sequenceId => {
            const sequence = Sequences.get(sequenceId);

            if (!sequence) return null;

            return (
                <SidebarSequenceTemplate key={sequence.Id} sequence={sequence} />
            );
        })}
        </div>
    );
};

interface SidebarActionTemplateProps {
    action: Action_f;
}

interface SidebarSequenceTemplateProps {
    sequence: Sequence_f;
}

interface SidebarApplicationTemplateProps {
    application: ChildApplicationInfo;
}

const SidebarActionTemplate = React.memo(({ action }: SidebarActionTemplateProps) => {
    return (
        <div>
        </div>
    );
});

const SidebarSequenceTemplate = React.memo(({ sequence }: SidebarSequenceTemplateProps) => {
    return (
        <div>
        </div>
    );
});

const SidebarApplicationTemplate = React.memo(({ application }: SidebarApplicationTemplateProps) => {
    return (
        <div>
        </div>
    );
});