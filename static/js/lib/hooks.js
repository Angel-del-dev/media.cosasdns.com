export const useState = initial_state => {
    let state = initial_state;

    return [
        () => state,
        new_state => {state = new_state;}
    ];
};