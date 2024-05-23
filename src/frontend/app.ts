// Matieral Web Components
import '@material/web/button/elevated-button';
import '@material/web/list/list';
import '@material/web/list/list-item';
import '@material/web/tabs/tabs';
import '@material/web/tabs/primary-tab';
import '@material/web/tabs/secondary-tab';
import '@material/web/textfield/outlined-text-field';
// Compile with: npx tsc src/frontend/app.ts && npx rollup -p @rollup/plugin-node-resolve ./dist/app.js -o ./dist/bundle.js
// View Functions
const accountView = (): void => {
    const account = document.createElement('account');
    account.innerHTML = `
        <div class="view-top-padding"></div>
        <div class="view-input-container">
            <md-outlined-text-field label="Username" type="text">
            </md-outlined-text-field>
            <md-outlined-text-field label="Password" type="password">
            </md-outlined-text-field>
        </div>
        <md-elevated-button type="submit">Login</md-elevated-button>
    `;
    account.classList.add('view');
    document.getElementById('app')!.appendChild(account);
};
const actionsView = (): void => {
    const actions = document.createElement('div');
    actions.innerHTML = `
        <md-tabs>
            <md-secondary-tab role="button" id="create">Create Record</md-secondary-tab>
            <md-secondary-tab role="button" id="modify">Modify Record</md-secondary-tab>
            <md-secondary-tab role="button" id="view">View Record</md-secondary-tab>
        </md-tabs>
        <div class="view-top-padding"></div>
    `;
    actions.classList.add('view');
    document.getElementById('app')!.appendChild(actions);
};
const supportView = (): void => {
    const support = document.createElement('div');
    support.innerHTML = `
        <div class="view-top-padding"></div>
        <md-list>
            <md-list-item>Phone: 555-555-5555</md-list-item>
            <md-list-item>Email: example@example.com</md-list-item>
            <md-list-item>Address: 1234 Example St.</md-list-item>
            <md-list-item>Hours: 9am-5pm</md-list-item>
        </md-list>
    `;
    support.classList.add('view');
    document.getElementById('app')!.appendChild(support);
};

// Call for views based on requested path
const showView = (hash: string): void => {
    // clear parent element
    document.getElementById('app')!.innerHTML = '';
    // choose new view based on hash
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
            document.getElementById('app')!.innerHTML = 'Home Page';
            break;
        default:
            document.getElementById('app')!.innerHTML = '404 Not Found';
    };
};

// Main
const app = (): void => {
    showView(window.location.hash);
    window.onhashchange = () => {
        showView(window.location.hash);
    };
};

app();