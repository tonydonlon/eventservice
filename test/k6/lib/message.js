
export const SESSION_START = 'SESSION_START';
export const SESSION_END = 'SESSION_END';
export const EVENT = 'EVENT';

const eventNames = [
    'event1',
    'event2',
    'test'
];

export function createMessageStream(sessionId, includeStart, includeEnd, numberOfEvents) {
    let messages = [];

    for (let i = 0; i < numberOfEvents; i++) {
        const eventName = eventNames[Math.floor(Math.random() * eventNames.length)];
        messages.push({
            time: Date.now(),
            type: EVENT,
            name: eventName,
        });
    }

    if (includeEnd) {
        messages.push({
            time: Date.now(),
            type: SESSION_END,
            session_id: sessionId,
        });
    }

    if (includeStart) {
        messages.unshift({
            time: Date.now(),
            type: SESSION_START,
            session_id: sessionId,
        });
    }

    return messages;
}
