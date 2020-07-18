import React from 'react';
import { withStyles } from 'material-ui/styles';
import Button from 'material-ui/Button';
import HouseService from '../../services/house.service';
import Paper from 'material-ui/Paper';
import Grid from 'material-ui/Grid';
import AppBar from 'material-ui/AppBar';
import Toolbar from 'material-ui/Toolbar';
import Typography from 'material-ui/Typography';
import orange from 'material-ui/colors/deepOrange';
import { Link } from 'react-router-dom';

const styles = theme => ({
    button: {
        margin: theme.spacing.unit,
    },
    flex: {
        flex: 1,
    },
    root: {
        margin: 'auto',
        backgroundSize: 'cover',
        backgroundRepeat: 'no-repeat',
        height: '70vh',
        padding: '320px 480px',
        backgroundImage: "url('https://cdn.freshome.com/wp-content/uploads/2016/10/white-houses-freshome18.png')",
        textAlign: "center",
    },
    paper1: {
        margin: 'auto',
        marginTop: theme.spacing.unit * 3,
        marginLeft: '12px',
        marginRight: '12px',
        padding: '160px 160px',
        height: '320px',
        backgroundImage: "url('https://www.careergirldaily.com/wp-content/uploads/2017/05/8e2e48fd3c6432d6a6a0b575fc14a269.jpg')",
        textAlign: 'center',
        opacity: '0.7',
    },
    paper2: {
        margin: 'auto',
        marginTop: theme.spacing.unit * 3,
        marginLeft: '12px',
        marginRight: '12px',
        padding: '160px 160px',
        height: '320px',
        backgroundImage: "url('https://img.domino.com/serve/the-best-kitchens-on-instagram-currently-568c09519ce40985052063f4-w620_h800.jpg')",
        textAlign: 'center',
        opacity: '0.7',
    },
    paper3: {
        margin: 'auto',
        marginTop: theme.spacing.unit * 3,
        marginLeft: '12px',
        marginRight: '12px',
        padding: '160px 160px',
        height: '320px',
        backgroundImage: "url('https://2.bp.blogspot.com/-e3mxtRr4dsY/WNzXPYcM5mI/AAAAAAABCFs/AbqMPNtGDf4I_kFXEExPPRK9RILgQJjdwCLcB/s1600/victoria6%25C2%25A9carinaolander.jpg')",
        textAlign: 'center',
        opacity: '0.7',
    },
    text: {
        color: orange[500],
    },
});

const Public = (props) => {
    const { classes } = props;
    const service = new HouseService();
    const showTupleCount = () => {
        service.getTotalTupleCount()
            .then(result => alert(
                'acc_user: ' + result.acc_user + '\n' + 
                'buyer: ' + result.buyer + '\n' + 
                'seller: ' + result.seller + '\n' + 
                'house: ' + result.house + '\n' + 
                'likes: ' + result.likes + '\n' + 
                'viewed: ' + result.viewed + '\n' + 
                'bought_house: ' + result.bought_house + '\n' + 
                'Total tuples: ' + result.total)    
            )
            .catch(err => alert('Something went wrong, please try again.'));
    }

    return ( 
        <div>
            <AppBar position="static" color="primary">
                <Toolbar>
                    <div className={classes.flex}>
                        <Typography variant="title" color="inherit">
                            HomeFolio
                        </Typography>
                    </div>
                    <Button
                        variant='raised'
                        secondary='true'
                        className={classes.button}
                        color='secondary'
                        size='small'
                        onClick={showTupleCount} >
                        Total Tuples
                    </Button>
                </Toolbar>
            </AppBar>

            <Paper className={classes.root} elevation={0}>
                <Typography variant="display3">
                    Welcome to HomeFolio!
                </Typography>
                <Typography variant="headline">
                    Please
                </Typography>
                <Button
                    variant='raised'
                    primary='true'
                    className={classes.button}
                    color='primary'
                    size='large'
                    component={Link}
                    to={'/login'} >
                    LOGIN
                </Button>
                <Button
                    variant='raised'
                    primary='true'
                    className={classes.button}
                    color='primary'
                    size='large'
                    component={Link}
                    to='/signup' >
                    SIGNUP
                </Button>
            </Paper>

            <Grid container>
                <Grid item xs={4}>
                    <Paper className={classes.paper1} elevation={0}>   
                        <Typography variant="display2" color="secondary">
                            130,000+ Houses 100,000+ Users
                        </Typography>
                    </Paper>
                </Grid>
                <Grid item xs={4}>
                    <Paper className={classes.paper2} elevation={0}>
                        <Typography variant="display2" color="primary">
                            Powerful Search and Compare
                        </Typography>
                    </Paper>
                </Grid>
                <Grid item xs={4}>
                    <Paper className={classes.paper3} elevation={0}>
                        <Typography variant="display2" className={classes.text}>
                            Great Seller Buyer Interaction
                        </Typography>
                    </Paper>
                </Grid>
            </Grid>

            <AppBar position="static" color="default">
                <Toolbar>
                    <Typography variant="body2" color="inherit">
                        © Group 25
                    </Typography>
                </Toolbar>
            </AppBar>
        </div>
    )
}

export default withStyles(styles)(Public);
