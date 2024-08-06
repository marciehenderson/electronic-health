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
const FORM_GENERATED_CONTAINER_SCRIPT_START = `\
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
`;
const FORM_GENERATED_CONTAINER_SCRIPT_END = `\
    xmlHttp.open('POST', '/action', true);
    xmlHttp.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xmlHttp.send('sub_hash=view&patient_id='+patient.value+'&record_date='+record.value);
}
`;
const FORM_GENERATED_CONTAINER_SCRIPT_VIEW = FORM_GENERATED_CONTAINER_SCRIPT_START + `\
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
` + FORM_GENERATED_CONTAINER_SCRIPT_END;
const FORM_GENERATED_CONTAINER_SCRIPT_MODIFY = FORM_GENERATED_CONTAINER_SCRIPT_START + `\
                <tr>\
                    <td>'+form[0].replace('patient_id=','')+'</td>\
                    <td>'+form[1].replace('record_date=','')+'</td>\
                    <td><textarea type=&quot;text&quot; class=&quot;generated-input&quot; name=&quot;location_id&quot; id=&quot;location_id&quot;></textarea></td>\
                    <td><textarea type=&quot;text&quot; class=&quot;generated-input&quot; name=&quot;record_type&quot; id=&quot;record_type&quot;></textarea></td>\
                    <td><textarea type=&quot;text&quot; class=&quot;generated-input&quot; name=&quot;notes&quot; id=&quot;notes&quot;></textarea></td>\
                </tr>\
            </table>';
            let generatedInputs = document.getElementsByClassName('generated-input');
            for (let i=0; i<generatedInputs.length; i++) {
                let input = generatedInputs[i];
                input.value = form[i+2].substring(form[i+2].indexOf('=')+1);
            }
        }
    }
` + FORM_GENERATED_CONTAINER_SCRIPT_END;
const accountView = () => {
    const account = document.createElement('account');
    account.innerHTML = `
        <div class="view-top-padding"></div>
        <form class="view-input-container" action="/login" method="post" onsubmit="
            async function fetchUserData(form) {
                console.log('Fetching user data...');
                const options = {
                    method: 'get',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Basic' + btoa(form.username.value + ':' + form.password.value)
                    }
                }
                await fetch('/userdata', options).then((response) => {
                    if (!response.ok) {
                        throw new Error('HTTP error: ' + response.status);
                    }
                        let json = response.json();
                        console.log('Response:', json);
                        return json;
                }).then((data) => {
                    // store user data with indexedDB
                    let request = indexedDB.open('user_data');
                    request.onupgradeneeded = function(event) {
                        let db = event.target.result;
                        let objectStore = db.createObjectStore('data', { keyPath: 'id' });
                        objectStore.createIndex('user_hash', 'user_hash', { unique: true });
                        objectStore.createIndex('password_hash', 'password_hash', { unique: false });
                    };
                    request.onerror = function(event) {
                        console.log('Database error: ' + event.target.errorCode);
                    };
                    request.onsuccess = function(event) {
                        let db = event.target.result;
                        let objectStore = db.transaction('data', 'readwrite').objectStore('data');
                        console.log('Adding data:', data);
                        let request = objectStore.add(JSON.parse(data));
                        request.onsuccess = function(event) {
                            console.log('Data added:', event.target.result);
                        };
                        request.onerror = function(event) {
                            console.log('Data error:', event.target.errorCode);
                        };
                        form.submit();
                    };
                }).catch((error) => {
                    console.error('Error:', error);
                });
            }
            fetchUserData(this);
            return false;
        ">
            <md-outlined-text-field name="username" label="Username" type="text" required>
            </md-outlined-text-field>
            <md-outlined-text-field name="password" label="Password" type="password" required>
            </md-outlined-text-field>
            <md-elevated-button type="submit">Login</md-elevated-button>
        </form>
    `;
    account.classList.add('view');
    document.getElementById('app').appendChild(account);
};
const actionsView = (subhash) => {
    const record = [generateOptions('record_date', 'record_data'), generateOptions('patient_id', 'record_data')];
    const patient = [generateOptions('patient_id', 'client_data'), generateOptions('last_name', 'patient_data'), generateOptions('first_name', 'patient_data')];
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
    var actionFormInner = `
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
    for (let i = 0; i < patient[0].length; i++) {
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
            actionFormInner = `
                <div class="action-indicator-bar">
                    <div></div>
                    <div class="action-indicator"></div>
                    <div></div>
                </div>
            ` + actionFormInner + `
                <input name="sub_hash" type="text" value="modify" style="display: none;"></input>
                <md-outlined-select name="record_date" label="Record Date" id="record_date" type="text" required onchange="
                    ${FORM_GENERATED_CONTAINER_SCRIPT_MODIFY}
                ">
                    <md-select-option value="-1" class="date-option" id="date-option-pid-default-">
                        <div slot="headline"></div>
                    </md-select-option>
            `;
            for (let i = 0; i < record[0].length; i++) {
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
            actionFormInner = `
                <div class="action-indicator-bar">
                    <div></div>
                    <div></div>
                    <div class="action-indicator"></div>
                </div>
            ` + actionFormInner + `
                <input name="sub_hash" type="text" value="view" style="display: none;"></input>
                <md-outlined-select name="record_date" label="Record Date" id="record_date" type="text" required onchange="
                    ${FORM_GENERATED_CONTAINER_SCRIPT_VIEW}
                ">
                    <md-select-option value="-1" class="date-option" id="date-option-pid-default-">
                        <div slot="headline"></div>
                    </md-select-option>
            `;
            for (let i = 0; i < record[0].length; i++) {
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
    }
    ;
    if (subhash !== 'view') {
        actionFormInner += `
            <md-elevated-button type="submit">Submit</md-elevated-button>
        </form>
    `;
    }
    else {
        actionFormInner += `
            </form>
        `;
    }
    actionsSubView.innerHTML = actionFormInner;
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
            <md-list-item>
                <md-filled-text-field label="Contact Support" type="textarea">
                </md-filled-text-field>
            </md-list-item>
            <md-elevated-button type="submit">Submit</md-elevated-button>
        </md-list>
    `;
    support.classList.add('view');
    document.getElementById('app').appendChild(support);
};
const setIndicator = (hash, id) => {
    var indicator = document.getElementById(id);
    var active = hash.substring(1).concat('-indicator');
    if (indicator.classList.contains('view-indicator')) {
        if (id !== active) {
            indicator.classList.remove('view-indicator');
        }
    }
    else if (id === active) {
        indicator.classList.add('view-indicator');
    }
};
const generateOptions = (column, cookie) => {
    const cookies = `; ${document.cookie};`;
    const dIndex = cookies.indexOf(`; ${cookie}=`) + 1;
    const data = cookies.substring(dIndex, cookies.indexOf(';', dIndex)).split('],[');
    let options = [];
    data.forEach((row) => {
        const cIndex = row.indexOf(`${column}=`);
        options.push(row.substring(cIndex, row.indexOf(',', cIndex)).substring(column.length + 1));
    });
    return options;
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
    setIndicator(hash, 'account-indicator');
    setIndicator(hash, 'actions-indicator');
    setIndicator(hash, 'support-indicator');
};
const app = () => {
    showView(window.location.hash);
    window.onhashchange = () => {
        showView(window.location.hash);
    };
};
app();
//# sourceMappingURL=app.js.map