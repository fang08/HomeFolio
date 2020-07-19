import React from 'react';
import { Provider } from "react-redux";
import store from "./redux/store/main";
import { BrowserRouter, Route, Switch, Redirect } from 'react-router-dom';
import Home from './scenes/Home/Home.component';
import HouseInfo from './scenes/HouseInfo/HouseInfo.component';
import Login from './scenes/Login/Login.component';
import Signup from './scenes/Signup/Signup.component';
import Public from './scenes/Public/Public.component';
import Compare from './scenes/Compare/Compare.component';
import Search from './scenes/Search/Search.component';
import UserInfo from './scenes/UserInfo/UserInfo.component';
import NotFound from './scenes/NotFound/NotFound.component';
import AuthService from './services/auth.service';

export const servAddr = 'http://ec2-3-128-153-42.us-east-2.compute.amazonaws.com:1323';
// export const servAddr = 'http://localhost:1323';
export const urlPrefix = '/api/v1';

class Main extends React.Component {
    constructor(props) {
        super(props);
        this.service = new AuthService();
    }

    PrivateRoute = ({ component: Component, ...rest }) => (
        <Route {...rest} render={props => (
            this.service.authenticate(localStorage.getItem('auth_token')) ? 
                (<Component {...props} />) : (<Redirect to={{ pathname: "/login" }} />)
        )} />
    );

    render() {
        return (
            <Provider store={store}>
                <BrowserRouter basename={process.env.PUBLIC_URL}>
                    <Switch>
                        <Route exact path='/' component={Public} />
                        <Route exact path='/login' component={Login} />
                        <Route exact path='/signup' component={Signup} />
                        <this.PrivateRoute path='/home' component={Home} />
                        <this.PrivateRoute path='/houseinfo/:h_id' component={HouseInfo} />
                        <this.PrivateRoute path='/userinfo/:u_id' component={UserInfo} location={this.props.location} />
                        <this.PrivateRoute path='/compare' component={Compare} />
                        <this.PrivateRoute path='/search' component={Search} />
                        <Route component={NotFound} />
                    </Switch>
                </BrowserRouter>
            </Provider>
        )
    }
}

export default Main;
