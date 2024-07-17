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
    // Set dropdown options based on database records
    const record = [generateOptions('record_date','record_data'), generateOptions('patient_id','record_data')];
    const patient = [generateOptions('patient_id','client_data'), generateOptions('last_name','patient_data'), generateOptions('first_name','patient_data')];
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
            <md-outlined-select name="patient_id" label="Patient ID" id="patient_id" type="text" required onchange="
                let patient = document.getElementById('patient_id').value;
                let records = document.getElementsByClassName('date-option');
                for (let i=0; i<records.length; i++) {
                    let record = records[i];
                    if (record.id.includes('pid-'+patient+'-')) {
                        record.style.display = 'block';
                    } else if (!record.id.includes('default')) {
                        record.style.display = 'none';
                    }
                }
                let recordDate = document.getElementById('record_date');
                recordDate.selectedIndex = 0;
            ">
                <md-icon slot="trailing-icon">search</md-icon>
    `;
    for (let i=0; i<patient[0].length; i++) {
        let value = patient[0][i] + ': ' + patient[1][i] + ', ' + patient[2][i];
        actionFormInner += `
            <md-select-option value="${patient[0][i]}">
                <div slot="headline">${value}</div>
            </md-select-option>
        `;
    }
    actionFormInner += `
            </md-outlined-select>
    `;
    switch (subhash) {
        default:
        case 'create':
            // prepend active indicator bar and append form fields
            actionFormInner = `
                <div class="action-indicator-bar">
                    <div class="action-indicator"></div>
                    <div></div>
                    <div></div>
                </div>
            ` + actionFormInner + `
                <input name="sub_hash" type="text" value="create" style="display: none;"></input>
                <md-outlined-text-field name="location_id" label="Location ID" type="text" required>
                    <md-icon slot="trailing-icon">search</md-icon>
                </md-outlined-text-field>
                <md-outlined-select name="record_type" label="Record Type" type="text" required>
                    <md-select-option value="check-up">
                        <div slot="headline">Check-Up</div>
                    </md-select-option>
                    <md-select-option value="annual">
                        <div slot="headline">Annual</div>
                    </md-select-option>
                    <md-select-option value="blood-work">
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
                    <div class="action-indicator"></div>
                    <div></div>
                </div>
            ` + actionFormInner + `
                <input name="sub_hash" type="text" value="modify" style="display: none;"></input>
                <md-outlined-select name="record_date" label="Record Date" id="record_date" type="text" required onchange="
                    let patient = document.getElementById('patient_id');
                    let record = document.getElementById('record_date');
                    if (patient.value !== '-1' && record.value !== '-1') {
                        let xmlHttp = new XMLHttpRequest();
                        xmlHttp.onreadystatechange = function() {
                            if (xmlHttp.readyState == 4 && xmlHttp.status == 200) {
                                const cookies = '; ' + document.cookie + ';';
                                const dIndex = cookies.indexOf('; record_view=') + 1;
                                const data = cookies.substring(dIndex, cookies.indexOf(';', dIndex)).split('],[');
                                let form = data[0].substring(15, data[0].length-3).split(',');
                                document.getElementById('form-generated-container').innerHTML = 
                                '<table>\
                                    <tr>\
                                        <th>Patient</th>\
                                        <th>Date</th>\
                                        <th>Location</th>\
                                        <th>Type</th>\
                                        <th>Notes</th>\
                                    </tr>\
                                    <tr>\
                                        <td>'+form[0].replace('patient_id=','')+'</td>\
                                        <td>'+form[1].replace('record_date=','')+'</td>\
                                        <td>'+form[2].replace('location_id=','')+'</td>\
                                        <td>'+form[3].replace('record_type=','')+'</td>\
                                        <td>'+form[4].replace('notes=','')+'</td>\
                                    </tr>\
                                </table>';
                            }
                        }
                        xmlHttp.open('POST', '/action', true);
                        xmlHttp.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
                        xmlHttp.send('sub_hash=view&patient_id='+patient.value+'&record_date='+record.value);
                    }
                ">
                    <md-select-option value="-1" class="date-option" id="date-option-pid-default-">
                        <div slot="headline"></div>
                    </md-select-option>
            `;
            for (let i=0; i<record[0].length; i++) {
                actionFormInner += `
                    <md-select-option value="${record[0][i]}" class="date-option" id="date-option-pid-${record[1][i]}-" style="display:none;">
                        <div slot="headline">${record[0][i]}</div>
                    </md-select-option>
                `;
            }
            actionFormInner += `
                </md-outlined-select>
                <div id="form-generated-container"></div>
            `;
            break;
        case 'view':
            // prepend active indicator bar and append form fields
            actionFormInner = `
                <div class="action-indicator-bar">
                    <div></div>
                    <div></div>
                    <div class="action-indicator"></div>
                </div>
            ` + actionFormInner + `
                <input name="sub_hash" type="text" value="view" style="display: none;"></input>
                <md-outlined-select name="record_date" label="Record Date" id="record_date" type="text" required>
                    <md-select-option value="-1" class="date-option" id="date-option-pid-default-">
                        <div slot="headline"></div>
                    </md-select-option>
            `;
            for (let i=0; i<record[0].length; i++) {
                actionFormInner += `
                    <md-select-option value="${record[0][i]}" class="date-option" id="date-option-pid-${record[1][i]}-" style="display:none;">
                        <div slot="headline">${record[0][i]}</div>
                    </md-select-option>
                `;
            }
            actionFormInner += `
                </md-outlined-select>
                <div id="form-generated-container"></div>
            `;
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
// Set active indicator bar based on hash
const setIndicator = (hash: string, id: string): void => {
    // get the appropriate indicator element
    var indicator: HTMLDivElement = document.getElementById(id) as HTMLDivElement;
    // get id of active view indicator
    var active: string = hash.substring(1).concat('-indicator');
    // if the indicator is already active, remove it if the id doesn't match
    if (indicator.classList.contains('view-indicator')) {
        if(id !== active) {
            indicator.classList.remove('view-indicator');
        }
    }
    // if the indicator is not active, add it if the id matches
    else if (id === active) {
        indicator.classList.add('view-indicator');
    }
}
// Generate dropdown options based on database records
const generateOptions = (column: string, cookie: string): string[] => {
    // fetch options from database
    const cookies = `; ${document.cookie};`;
    // get all data from specified cookie
    const dIndex = cookies.indexOf(`; ${cookie}=`) + 1;
    // store each row of data as an element in an array
    const data = cookies.substring(dIndex, cookies.indexOf(';', dIndex)).split('],[');
    // get all instances of the specified column
    let options: string[] = [];
    data.forEach((row: string) => {
        const cIndex = row.indexOf(`${column}=`);
        options.push(row.substring(cIndex, row.indexOf(',', cIndex)).substring(column.length + 1));
    });
    return options;
}
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
    // set view indicator bar based on hash
    setIndicator(hash, 'account-indicator');
    setIndicator(hash, 'actions-indicator');
    setIndicator(hash, 'support-indicator');
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