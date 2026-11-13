import { useGlobalStore } from "../../../global";
import { TabsClasses } from "../Tester";
import "./Tabs.css"


interface TabsProps {
}

const Tabs = ({}:TabsProps) => {
    const tabs = useGlobalStore(state => state.Tabs);
    const setCurrentTab = useGlobalStore(state => state.setCurrentTab);
    const removeFromTabs = useGlobalStore(state => state.removeFromTabs);

    return (
        <div className={TabsClasses.container}>
            {tabs.map(tab => (
                <div key={tab.TypeId} onClick={() => setCurrentTab(tab)}>
                    {tab.TabType} {tab.Name}
                    <button onClick={(e) => {
                        e.stopPropagation();
                        removeFromTabs(tab);
                    }}>X</button>
                </div>
            ))}
        </div>
    );
};

export default Tabs;