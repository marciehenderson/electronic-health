// Matieral Web Components
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
// Compile with: npx tsc src/frontend/app.ts && npx rollup -p @rollup/plugin-node-resolve ./dist/app.js -o ./dist/bundle.js
// View Functions
// Account View - User Login Page
const accountView = (): void => {
    const account = document.createElement('account');
    account.innerHTML = `
        <div class="view-top-padding"></div>
        <form class="view-input-container" action="/login" method="post">
            <md-outlined-text-field name="username" label="Username" type="text" required>
            </md-outlined-text-field>
            <md-outlined-text-field name="password" label="Password" type="password" required>
            </md-outlined-text-field>
            <md-elevated-button type="submit">Login</md-elevated-button>
        </form>
    `;
    account.classList.add('view');
    document.getElementById('app')!.appendChild(account);
};
// Actions View - Create, Modify, View Records
const actionsView = (subhash: string): void => {
    const actions = document.createElement('div');
    actions.innerHTML = `
        <md-tabs>
            <md-secondary-tab role="button" onclick="window.location.hash='actions+create'" id="create">Create</md-secondary-tab>
            <md-secondary-tab role="button" onclick="window.location.hash='actions+modify'" id="modify">Modify</md-secondary-tab>
            <md-secondary-tab role="button" onclick="window.location.hash='actions+view'" id="view">View</md-secondary-tab>
        </md-tabs>
    `;
    actions.classList.add('view');
    document.getElementById('app')!.appendChild(actions);
    const actionsSubView = document.createElement('div');
    var actionFormInner: string = `
        <div class="view-top-padding"></div>
        <form class="view-input-container" action="/action" method="post">
            <md-outlined-text-field name="patient_id" label="Patient ID" type="text" required>
                <md-icon slot="trailing-icon">search</md-icon>
            </md-outlined-text-field>
    `;
    switch (subhash) {
        default:
        case 'create':
            // prepend active indicator bar and append form fields
            actionFormInner = `
                <div class="action-indicator-bar">
                    <div class="active-indicator"></div>
                    <div></div>
                    <div></div>
                </div>
            ` + actionFormInner + `
                <input name="sub_hash" type="text" value="create" style="display: none;"></input>
                <md-outlined-text-field name="location_id" label="Location ID" type="text" required>
                    <md-icon slot="trailing-icon">search</md-icon>
                </md-outlined-text-field>
                <md-outlined-select name="record_type" label="Record Type" type="text" required>
                    <md-select-option value="checkup">
                        <div slot="headline">Check-Up</div>
                    </md-select-option>
                    <md-select-option value="annual">
                        <div slot="headline">Annual</div>
                    </md-select-option>
                    <md-select-option value="bloodwork">
                        <div slot="headline">Blood-Work</div>
                    </md-select-option>
                    <md-select-option value="vaccination">
                        <div slot="headline">Vaccination</div>
                    </md-select-option>
                    <md-select-option value="emergency">
                        <div slot="headline">Emergency</div>
                    </md-select-option>
                </md-outlined-select>
                <md-outlined-text-field name="notes" label="Notes" type="textarea">
                </md-outlined-text-field>
            `;
            break;
        case 'modify':
            // prepend active indicator bar and append form fields
            actionFormInner = `
                <div class="action-indicator-bar">
                    <div></div>
                    <div class="active-indicator"></div>
                    <div></div>
                </div>
            ` + actionFormInner + `
                <input name="sub_hash" type="text" value="modify" style="display: none;"></input>
                <md-outlined-select name="record_date" label="Record Date" type="text" required>
                    <md-select-option value="1">
                        <div slot="headline">01/01/2021</div>
                    </md-select-option>
                    <md-select-option value="2">
                        <div slot="headline">01/02/2021</div>
                    </md-select-option>
                    <md-select-option value="3">
                        <div slot="headline">01/03/2021</div>
                    </md-select-option>
                    <md-select-option value="4">
                        <div slot="headline">01/04/2021</div>
                    </md-select-option>
                    <md-select-option value="5">
                        <div slot="headline">01/05/2021</div>
                    </md-select-option>
                </md-outlined-select>
                <md-outlined-select name="edit_value" label="Edit Value" type="text" required>
                    <md-select-option value="1">
                        <div slot="headline">Location ID</div>
                    </md-select-option>
                    <md-select-option value="2">
                        <div slot="headline">Record Type</div>
                    </md-select-option>
                    <md-select-option value="3">
                        <div slot="headline">Notes</div>
                    </md-select-option>
                </md-outlined-select>
                <div id="form-generated-container"></div>
            `;
            // note dates are hardcoded for example, should 
            // be generated from existing database records.
            // the date and edit value should be unselectable
            // until a patient ID is entered.
            // the form-generated-container will be populated
            // with the appropriate form fields based on the
            // selected edit value.
            break;
        case 'view':
            // prepend active indicator bar and append form fields
            actionFormInner = `
                <div class="action-indicator-bar">
                    <div></div>
                    <div></div>
                    <div class="active-indicator"></div>
                </div>
            ` + actionFormInner + `
                <input name="sub_hash" type="text" value="view" style="display: none;"></input>
                <md-outlined-select name="record_date" label="Record Date" type="text" required>
                    <md-select-option value="1">
                        <div slot="headline">01/01/2021</div>
                    </md-select-option>
                    <md-select-option value="2">
                        <div slot="headline">01/02/2021</div>
                    </md-select-option>
                    <md-select-option value="3">
                        <div slot="headline">01/03/2021</div>
                    </md-select-option>
                    <md-select-option value="4">
                        <div slot="headline">01/04/2021</div>
                    </md-select-option>
                    <md-select-option value="5">
                        <div slot="headline">01/05/2021</div>
                    </md-select-option>
                </md-outlined-select>
                <div id="form-generated-container"></div>
            `;
            // the date should be unselectable until a patient
            // ID is entered.
            // the form-generated-container will be populated
            // with the appropriate form fields based on the
            // selected record.
            break;
    };
    // append submit button and closing tag
    actionFormInner += `
            <md-elevated-button type="submit">Submit</md-elevated-button>
        </form>
    `;
    actionsSubView.innerHTML = actionFormInner;
    actionsSubView.classList.add('subview');
    document.getElementsByClassName('view')[0]!.appendChild(actionsSubView);
};
// Support View - Contact Information and User Support
const supportView = (): void => {
    const support = document.createElement('div');
    support.innerHTML = `
        <div class="view-top-padding"></div>
        <md-list>
            <md-list-item>Phone: 555-555-5555</md-list-item>
            <md-list-item>Email: example@example.com</md-list-item>
            <md-list-item>Address: 1234 Example St.</md-list-item>
            <md-list-item>Hours: 9am-5pm</md-list-item>
            <md-list-item>
                <md-filled-text-field label="Contact Support" type="textarea">
                </md-filled-text-field>
            </md-list-item>
            <md-elevated-button type="submit">Submit</md-elevated-button>
        </md-list>
    `;
    support.classList.add('view');
    document.getElementById('app')!.appendChild(support);
};
// Call for views based on requested path
const showView = (hash: string): void => {
    // clear parent element
    document.getElementById('app')!.innerHTML = '';
    // separate hash from sub-hash
    var subhash = hash.split('+')[1];
    hash = hash.split('+')[0];
    // choose new view based on hash
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
            document.getElementById('app')!.innerHTML = '404 Not Found';
    };
};
// Event Listener Functions
const activeTab = (tab: HTMLElement): void => {
    console.log(tab.id)
    // Add active attribute to active tab which doesn't have it
    if(window.location.hash.split('#')[1] === tab.id) {
        console.log('active')
        if(!tab.hasAttribute('active')) {
            tab.setAttribute('active', '');
            console.log(tab.attributes)
        }
    }
    // Remove active attribute from non-active tab that has it
    else if(tab.hasAttribute('active')) {
        tab.removeAttribute('active');
        console.log('inactive')
        console.log(tab.attributes)
    }
};
// Main
const app = (): void => {
    showView(window.location.hash);
    window.onhashchange = () => {
        showView(window.location.hash);
    };
};
// Run App
app();