import React, { Component } from 'react';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import RaisedButton from 'material-ui/TextField';
import TextField from 'material-ui/TextField/TextField';

class LoginSignup extends Component {
    render() {
        return (
            <div>
                <MuiThemeProvider>
                    <TextField placeholder='username' />
                </MuiThemeProvider>
            </div>
        )
    }
}

export default LoginSignup;