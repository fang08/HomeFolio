import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import TextField from 'material-ui/TextField';
import Button from 'material-ui/Button';
import AppBar from 'material-ui/AppBar';
import Typography from 'material-ui/Typography';
import Toolbar from 'material-ui/Toolbar';
import AuthService from '../../services/auth.service';
import { Redirect } from 'react-router-dom';

const styles = theme => ({
    container: {
        display: 'flex',
        flexWrap: 'wrap',
        textAlign: 'center',
    },
    textField: {
        marginLeft: theme.spacing.unit,
        marginRight: theme.spacing.unit,
        width: 200,
    },
    button: {
        margin: theme.spacing.unit,
    },
    root: {
        flexGrow: 1,
    },
    flex: {
        flex: 1,
    }
});

class Login extends React.Component {
    constructor(props) {
        super(props);
        this.service = new AuthService();
        this.state = {
            redirectToHome: false,
            email: '',
            password: ''
        }
    }

    login = () => {
        this.service.login(this.state.email, this.state.password)
            .then(() => this.setState({ redirectToHome: true }))
            .catch((err) => alert(err.message))
    }

    render() {
        const { classes } = this.props;

        if (this.state.redirectToHome)
            return <Redirect to={{ pathname: "/home" }} />

        return (
            <div>
                <div className={classes.root} onKeyDown={e => {if(e.keyCode === 13) this.login();}}>
                    <AppBar position='static'>
                        <Toolbar>
                            <Typography variant="title" color="inherit" className={classes.flex}>
                                Login
                            </Typography>
                        </Toolbar>
                    </AppBar>
                    <form className={classes.container}>
                        <TextField
                            required
                            autoFocus
                            id='loginEmailField'
                            className={classes.textField}
                            label='email'
                            placeholder='Email' 
                            margin='normal'
                            onChange={event => this.setState({email: event.target.value})}
                        />
                        <TextField
                            required
                            id='loginPasswordField'
                            className={classes.textField}
                            label='password'
                            placeholder='Password'
                            type='password'
                            margin='normal'
                            onChange={event => this.setState({password: event.target.value})}
                        />
                        <Button
                            variant='raised'
                            primary='true'
                            className={classes.button}
                            color='primary'
                            onClick={this.login} >
                            login
                        </Button>
                    </form>
                </div>
                <div>
                    <br/>
                    Demo user:
                    <br/>
                    Email: test@example.com
                    <br/>
                    Password: test
                </div>
            </div>
        )
    }
}

Login.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(Login);
