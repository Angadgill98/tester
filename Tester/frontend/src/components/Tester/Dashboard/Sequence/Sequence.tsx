import React, { useEffect, useState } from 'react';
import { GetSequence } from "../../../../../wailsjs/go/main/App";
import { Tab, Sequence as SequenceType } from '../../../../models/application';

interface SequenceProps {
    tab: Tab;
}

const Sequence = ({ tab }: SequenceProps) => {
    const [sequence, setSequence] = useState<SequenceType | null>(null);

    useEffect(() => {
        let cancelled = false;

        setSequence(null);

        const fetchSequence = async () => {
            try {
                const response = await GetSequence(tab.TypeId);

                if (cancelled) return;

                if (!response.Success) {
                    console.error(response.Error);
                    return;
                }

                setSequence(response.Data);
            } catch (error) {
                if (!cancelled) {
                    console.error("Failed to fetch sequence:", error);
                }
            }
        };

        fetchSequence();

        return () => {
            cancelled = true;
        };
    }, [tab.TypeId, tab.AppId, tab.TabType]);

    if (!sequence) return <div>Loading sequence...</div>;

    return (
        <div>
            <h2>{sequence.SequenceName}</h2>
            <p>Sequence ID: {sequence.SequenceID}</p>
            <p>App ID: {sequence.AppID}</p>

            <h3>Actions</h3>
            {sequence.ActionIDs.map(actionId => (
                <div key={actionId}>
                    Action ID: {actionId}
                </div>
            ))}
        </div>
    );
};

export default Sequence;