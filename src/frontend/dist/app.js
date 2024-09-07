var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
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
    let fetchRecord = async (patient_id, record_date) => {
        let sub_hash = 'view';
        let formData = new FormData();
        formData.append('sub_hash', sub_hash);
        formData.append('patient_id', patient_id);
        formData.append('record_date', record_date);
        await fetch('/action', {
            method: 'post',
            headers: {
                'Content-Type': 'multipart/form-data'
            },
            body: formData
        }).then((response) => {
            if (!response.ok) {
                throw new Error('HTTP error: ' + response.status);
            }
            console.log('Response:', response);
            let json = response.json();
            console.log('JSON:', json);
            return json;
        }).then((json) => {
            let data = String(json).split(',');
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
const FORM_GENERATED_CONTAINER_SCRIPT_END = `
        }).catch((error) => {
            console.error('Error:', error);
        });
    }
    fetchRecord(patient.value, record.value);
}`;
const FORM_GENERATED_CONTAINER_SCRIPT_VIEW = FORM_GENERATED_CONTAINER_SCRIPT_START + `\
        <tr>\
            <td>'+data[0].substring(data[0].indexOf(':')+1).replace(new RegExp('&quot;', 'g'), '')+'</td>\
            <td>'+data[1].substring(data[1].indexOf(':')+1).replace(new RegExp('&quot;', 'g'), '')+'</td>\
            <td>'+data[2].substring(data[2].indexOf(':')+1).replace(new RegExp('&quot;', 'g'), '')+'</td>\
            <td>'+data[3].substring(data[3].indexOf(':')+1).replace(new RegExp('&quot;', 'g'), '')+'</td>\
            <td>'+data[4].substring(data[4].indexOf(':')+1).replace(new RegExp('&quot;', 'g'), '').replace('}', '')+'</td>\
        </tr>\
    </table>';
` + FORM_GENERATED_CONTAINER_SCRIPT_END;
const FORM_GENERATED_CONTAINER_SCRIPT_MODIFY = FORM_GENERATED_CONTAINER_SCRIPT_START + `\
        <tr>\
            <td>'+data[0].substring(data[0].indexOf(':')+1).replace(new RegExp('&quot;', 'g'), '')+'</td>\
            <td>'+data[1].substring(data[1].indexOf(':')+1).replace(new RegExp('&quot;', 'g'), '')+'</td>\
            <td><textarea type=&quot;text&quot; class=&quot;generated-input&quot; name=&quot;location_id&quot; id=&quot;location_id&quot;></textarea></td>\
            <td><textarea type=&quot;text&quot; class=&quot;generated-input&quot; name=&quot;record_type&quot; id=&quot;record_type&quot;></textarea></td>\
            <td><textarea type=&quot;text&quot; class=&quot;generated-input&quot; name=&quot;notes&quot; id=&quot;notes&quot;></textarea></td>\
        </tr>\
    </table>';
    let generatedInputs = document.getElementsByClassName('generated-input');
    for (let i=0; i<generatedInputs.length; i++) {
        let input = generatedInputs[i];
        input.value = data[i+2].substring(data[i+2].indexOf(':')+1).replace(new RegExp('&quot;', 'g'), '').replace('}', '');
    }
` + FORM_GENERATED_CONTAINER_SCRIPT_END;
const accountView = () => {
    const account = document.createElement('account');
    account.innerHTML = `
        <div class="view-top-padding"></div>
        <form class="view-input-container" action="/login" method="post" onsubmit="
            // asynchronous function to fetch user data and store in indexedDB
            async function fetchUserData(req, form, ver, store, key, index, next) {
                const options = {
                    method: 'get',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Basic' + btoa(form.username.value + ':' + form.password.value)
                    }
                }
                // return true if successful, false if not
                return await fetch(req, options).then((response) => {
                    if (!response.ok) {
                        throw new Error('HTTP error: ' + response.status);
                    }
                    let json = response.json();
                    return json;
                }).then((data) => {
                    // store user data with indexedDB
                    let request = indexedDB.open('user_data', ver);
                    request.onupgradeneeded = function(event) {
                        let db = event.target.result;
                        let objectStore = db.createObjectStore(store, { keyPath: key });
                        for (let i=0; i<index.length; i++) {
                            objectStore.createIndex(index[i][0], index[i][0], { unique: index[i][1] }); 
                        }
                    };
                    request.onerror = function(event) {
                        console.log('Database error: ' + event.target.errorCode);
                    };
                    request.onsuccess = function(event) {
                        let db = event.target.result;
                        let objectStore = db.transaction(store, 'readwrite').objectStore(store);
                        // if there is more than one row of data, split and add each row individually
                        if (data.includes('},{')) {
                            data = data.split('},');
                            for (let i=0; i<data.length; i++) {
                                if (!data[i].includes('}')) {
                                    data[i] = data[i] + '}';
                                }
                                let request = objectStore.add(JSON.parse(data[i]));
                                // if this is the last row of data, call next
                                if (i === data.length-1) {
                                    request.onsuccess = function(event) {
                                        next();
                                    };
                                } else {
                                    request.onsuccess = function(event) {
                                        // console.log('Data added:', event.target.result);
                                    };
                                }
                                request.onerror = function(event) {
                                    console.log('Data error:', event.target.errorCode);
                                };
                            }
                        }
                        // if there is only one row of data, add it
                        else {
                            let request = objectStore.add(JSON.parse(data));
                            request.onsuccess = function(event) {
                                next();
                            };
                            request.onerror = function(event) {
                                console.log('Data error:', event.target.errorCode);
                            };
                        }
                    };
                }).then(() => {
                    // return true if successful
                    return true;
                }).catch((error) => {
                    console.error('Error:', error);
                    return false;
                });
            }
            // call fetchUserData for each object store in sequence
            fetchUserData('/userdata', this, 1, 'credential', 'id', [['user_hash', true], ['password_hash', false]], () => {
            fetchUserData('/recorddata', this, 2, 'record', 'record_date', [['patient_id', false]], () => {
            fetchUserData('/clientdata', this, 3, 'client', 'patient_id', [['practitioner_id', false]], () => {
            fetchUserData('/patientdata', this, 4, 'patient', 'id', [['first_name', false], ['last_name', false], ['date_of_birth', false], ['street_address', false], ['contact_number', false], ['email', false]], () => {
            this.submit();});});});});
            // prevent default form submission
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
function actionsView(subhash) {
    return __awaiter(this, void 0, void 0, function* () {
        let record = yield Promise.all([generateOptions('record_date', 'record'), generateOptions('patient_id', 'record')]);
        let patient = yield Promise.all([generateOptions('patient_id', 'client'), generateOptions('last_name', 'patient'), generateOptions('first_name', 'patient')]);
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
        <form class="view-input-container" action="/action" method="post" onsubmit="
            if (document.getElementById('practitioner_id')) {
                let getUserID = (form) => {
                    // get value of practioner_id from indexedDB
                    db = indexedDB.open('user_data');
                    db.onsuccess = function(event) {
                        let db = event.target.result;
                        let objectStore = db.transaction('client', 'readonly').objectStore('client');
                        let request = objectStore.get(document.getElementById('patient_id').value);
                        request.onsuccess = function(event) {
                            let practitioner = event.target.result;
                            document.getElementById('practitioner_id').value = practitioner.practitioner_id;
                            form.submit();
                        };
                        request.onerror = function(event) {
                            console.log('Database error: ' + event.target.errorCode);
                        };
                    };
                }
                getUserID(this);
                return false;
            }
        ">
            <md-outlined-select name="patient_id" label="Patient ID" id="patient_id" type="text" required onchange="
                let patient = document.getElementById('patient_id').value;
                // only if the record_date element exists
                if (document.getElementById('record_date')) {
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
                }
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
                <input name="sub_hash" id="sub_hash" type="text" value="create" style="display: none;"></input>
                <input name="practitioner_id" id="practitioner_id" type="text" value="" style="display: none;"></input>
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
                <input name="sub_hash" id="sub_hash" type="text" value="modify" style="display: none;"></input>
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
                <input name="sub_hash" id="sub_hash" type="text" value="view" style="display: none;"></input>
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
    });
}
;
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
function generateOptions(column, store) {
    return __awaiter(this, void 0, void 0, function* () {
        let db = indexedDB.open('user_data');
        let promise = new Promise((resolve, reject) => {
            let options = [];
            db.onsuccess = function (event) {
                let db = event.target.result;
                let objectStore = db.transaction(store, 'readonly').objectStore(store);
                let request = objectStore.openCursor(null, 'nextunique');
                let row;
                request.onsuccess = function (event) {
                    let cursor = event.target.result;
                    if (cursor) {
                        row = JSON.stringify(cursor.value);
                        const cIndex = row.indexOf(`\"${column}\":\"`);
                        const vIndex = row.indexOf('\"', cIndex + column.length + 4);
                        const rowVal = row.substring(cIndex, vIndex).substring(column.length + 4);
                        options.push(rowVal);
                        cursor.continue();
                    }
                    else {
                        resolve(options);
                    }
                };
                request.onerror = function (event) {
                    console.log('Database error: ' + event.target.errorCode);
                    reject([]);
                };
            };
            db.onerror = function (event) {
                console.log('Database error: ' + event.target.errorCode);
                reject([]);
            };
        });
        return promise;
    });
}
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