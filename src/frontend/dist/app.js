"use strict";
const accountView = () => {
    const account = document.createElement('account');
    account.innerHTML = `
        <h2>Login</h2>
        <div>
            <label for="username">Username:</label>
            <input type="text" id="username" name="username">
        </div>
        <div>
            <label for="password">Password:</label>
            <input type="password" id="password" name="password">
        </div>
        <button type="submit">Login</button>
    `;
    document.getElementById('app').appendChild(account);
};
const actionsView = () => {
    const actions = document.createElement('div');
    actions.innerHTML = `
        <h2>Actions</h2>
        <button id="create">Create Record</button>
        <button id="modify">Modify Record</button>
        <button id="view">View Record</button>
    `;
    document.getElementById('app').appendChild(actions);
};
const supportView = () => {
    const support = document.createElement('div');
    support.innerHTML = `
        <h2>Support</h2>
        <ul>
            <li>Phone: 555-555-5555</li>
            <li>Email: example@example.com</li>
            <li>Address: 1234 Example St.</li>
            <li>Hours: 9am-5pm</li>
        </ul>
    `;
    document.getElementById('app').appendChild(support);
};
const showView = (hash) => {
    document.getElementById('app').innerHTML = '';
    switch (hash) {
        case '#account':
            accountView();
            break;
        case '#actions':
            actionsView();
            break;
        case '#support':
            supportView();
            break;
        case '':
            document.getElementById('app').innerHTML = 'Home Page';
            break;
        default:
            document.getElementById('app').innerHTML = '404 Not Found';
    }
    ;
};
const app = () => {
    showView(window.location.hash);
    window.onhashchange = () => {
        showView(window.location.hash);
    };
};
app();
//# sourceMappingURL=app.js.map