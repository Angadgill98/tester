import { create } from "zustand";
import { Action_f, Application_f, ChildApplicationInfo, Environment_f, Sequence_f } from "./models/application";




interface GlobalStore {
    Applications: Map<number, Application_f>;
    Actions: Map<number, Action_f>;
    Sequences: Map<number, Sequence_f>;
    Environments: Map<number, Environment_f>;
    ChildApplications: Map<number, ChildApplicationInfo>;

    GlobalApplication: Application_f | null;
    CurrentApplication: Application_f | null;

    setApplications: (id: number, app: Application_f) => void;
    setActions: (id: number, action: Action_f) => void;
    setSequences: (id: number, sequence: Sequence_f) => void;
    setEnvironments: (id: number, environment: Environment_f) => void;
    setChildApplications: (id: number, child: ChildApplicationInfo) => void;

    setGlobalApplication: (app: Application_f | null) => void;
    setCurrentApplication: (app: Application_f | null) => void;
}

export const useGlobalStore = create<GlobalStore>((set) => ({
    Applications: new Map<number, Application_f>(),
    Actions: new Map<number, Action_f>(),
    Sequences: new Map<number, Sequence_f>(),
    Environments: new Map<number, Environment_f>(),
    ChildApplications: new Map<number, ChildApplicationInfo>(),

    GlobalApplication: null,
    CurrentApplication: null,

    setApplications: (id, app) => set((state) => {
        const applications = new Map(state.Applications);
        applications.set(id, app);
        return { Applications: applications };
    }),

    setActions: (id, action) => set((state) => {
        const actions = new Map(state.Actions);
        actions.set(id, action);
        return { Actions: actions };
    }),

    setSequences: (id, sequence) => set((state) => {
        const sequences = new Map(state.Sequences);
        sequences.set(id, sequence);
        return { Sequences: sequences };
    }),

    setEnvironments: (id, environment) => set((state) => {
        const environments = new Map(state.Environments);
        environments.set(id, environment);
        return { Environments: environments };
    }),

    setChildApplications: (id, child) => set((state) => {
        const childApplications = new Map(state.ChildApplications);
        childApplications.set(id, child);
        return { ChildApplications: childApplications };
    }),

    setGlobalApplication: (app) => set({ GlobalApplication: app }),
    setCurrentApplication: (app) => set({ CurrentApplication: app })
}));