import '@material/web/icon/icon';
import '@material/web/button/elevated-button';
import '@material/web/list/list';
import '@material/web/list/list-item';
import '@material/web/tabs/tabs';
import '@material/web/tabs/primary-tab';
import '@material/web/tabs/secondary-tab';
import '@material/web/textfield/outlined-text-field';
import '@material/web/textfield/filled-text-field';
import '@material/web/select/outlined-select';
import '@material/web/select/select-option';
const accountView = () => {
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
    document.getElementById('app').appendChild(account);
};
const actionsView = (subhash) => {
    const actions = document.createElement('div');
    actions.innerHTML = `
        <md-tabs>
            <md-secondary-tab role="button" onclick="window.location.hash='actions+create'" id="create">Create</md-secondary-tab>
            <md-secondary-tab role="button" onclick="window.location.hash='actions+modify'" id="modify">Modify</md-secondary-tab>
            <md-secondary-tab role="button" onclick="window.location.hash='actions+view'" id="view">View</md-secondary-tab>
        </md-tabs>
    `;
    actions.classList.add('view');
    document.getElementById('app').appendChild(actions);
    const actionsSubView = document.createElement('div');
    switch (subhash) {
        default:
        case 'create':
            actionsSubView.innerHTML = `
                <div class="action-indicator-bar">
                    <div class="active-indicator"></div>
                    <div></div>
                    <div></div>
                </div>
                <div class="view-top-padding"></div>
                <div class="view-input-container">
                    <md-outlined-text-field label="Patient ID" type="search">
                        <md-icon slot="trailing-icon">search</md-icon>
                    </md-outlined-text-field>
                    <md-outlined-text-field label="Location ID" type="search">
                        <md-icon slot="trailing-icon">search</md-icon>
                    </md-outlined-text-field>
                    <md-outlined-select label="Type">
                        <md-select-option value="1">
                            <div slot="headline">Check-Up</div>
                        </md-select-option>
                        <md-select-option value="2">
                            <div slot="headline">Annual</div>
                        </md-select-option>
                        <md-select-option value="3">
                            <div slot="headline">Blood-Work</div>
                        </md-select-option>
                        <md-select-option value="4">
                            <div slot="headline">Vaccination</div>
                        </md-select-option>
                        <md-select-option value="5">
                            <div slot="headline">Emergency</div>
                        </md-select-option>
                    </md-outlined-select>
                    <md-outlined-text-field label="Notes" type="textarea">
                    </md-outlined-text-field>
                    <md-elevated-button type="submit">Submit</md-elevated-button>
                </div>
            `;
            break;
        case 'modify':
            actionsSubView.innerHTML = `
                <div class="action-indicator-bar">
                    <div></div>
                    <div class="active-indicator"></div>
                    <div></div>
                </div>
                <div class="view-top-padding"></div>
                <h1>Modify Record</h1>
            `;
            break;
        case 'view':
            actionsSubView.innerHTML = `
                <div class="action-indicator-bar">
                    <div></div>
                    <div></div>
                    <div class="active-indicator"></div>
                </div>
                <div class="view-top-padding"></div>
                <h1>View Record</h1>
            `;
            break;
    }
    ;
    actionsSubView.classList.add('subview');
    document.getElementsByClassName('view')[0].appendChild(actionsSubView);
};
const supportView = () => {
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
    document.getElementById('app').appendChild(support);
};
const showView = (hash) => {
    document.getElementById('app').innerHTML = '';
    var subhash = hash.split('+')[1];
    hash = hash.split('+')[0];
    switch (hash) {
        case '#account':
            accountView();
            break;
        case '#actions':
            actionsView(subhash);
            break;
        case '#support':
            supportView();
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