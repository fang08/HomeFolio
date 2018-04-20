import React from 'react';
import PropTypes from 'prop-types';
import { withStyles } from 'material-ui/styles';
import NavBar from '../../components/NavBar/NavBar.component';
import HouseService from '../../services/house.service';
import UserService from '../../services/user.service';
import classnames from 'classnames';
import { compose, withProps } from "recompose";
import { withScriptjs, withGoogleMap, GoogleMap, Marker, StreetViewPanorama } from "react-google-maps";
import Card, { CardHeader, CardContent, CardActions } from 'material-ui/Card';
import Collapse from 'material-ui/transitions/Collapse';
import Avatar from 'material-ui/Avatar';
import IconButton from 'material-ui/IconButton';
import Typography from 'material-ui/Typography';
import red from 'material-ui/colors/red';
import FavoriteIcon from '@material-ui/icons/Favorite';
import VisibilityIcon from '@material-ui/icons/Visibility';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import Button from 'material-ui/Button';
import AddIcon from '@material-ui/icons/Add';
import EditIcon from '@material-ui/icons/Edit'
import DeleteIcon from '@material-ui/icons/Delete';
import Tooltip from 'material-ui/Tooltip';
import Grid from 'material-ui/Grid';

const styles = theme => ({
    button: {
        margin: theme.spacing.unit,
    },
    root: {
        flexGrow: 1,
    },
    flex: {
        flex: 1,
    },
    card: {
        width: '50%',
        margin: '0 auto',
    },
    media: {
        height: 294,
    },
    actions: {
        display: 'flex',
    },
    expand: {
        transform: 'rotate(0deg)',
        transition: theme.transitions.create('transform', {
            duration: theme.transitions.duration.shortest,
        }),
        marginLeft: 'auto',
    },
    expandOpen: {
        transform: 'rotate(180deg)',
    },
    avatar: {
        backgroundColor: red[500],
    },

    fab: {
        margin: theme.spacing.unit * 2,
    },
    absolute: {
        position: 'absolute',
        bottom: theme.spacing.unit * 2,
        right: theme.spacing.unit * 3,
    },
});

class HouseInfo extends React.Component {
    constructor(props) {
        super(props);
        this.service = new HouseService();
        this.userService = new UserService();
        this.state = {
            info: {
                h_id: '',
                u_id: '',
                bathroomCnt: 0, 
                bedroomCnt: 0,
                buildingQualityID: 0,
                livingAreaSize: 0,
                latitude: 0.0,
                longitude: 0.0,
                lotSize: 0,
                cityID: 0,
                county: '',
                zip: 0,
                yearBuilt: 0,
                storyNum: 0,
                price: 0,
                tax: 0
            },
            userinfo: {
                u_id: '',
                email: '',
                username: '',
                password: '',
                age: 0,
                area: '',
                bio: '',
                seller: '',
                buyer: ''
            },
            addr: ''
        }
    }

    StreetViewPanorama = compose(
        withProps({
            googleMapURL: "https://maps.googleapis.com/maps/api/js?&key=AIzaSyAHbTvrtAr7iIMx0ZHhwwB3RqgWpRy4fvs",
            loadingElement: <div style={{ height: `100%` }} />,
            containerElement: <div style={{ height: `400px` }} />,
            mapElement: <div style={{ height: `100%` }} />,
        }),
            withScriptjs,
            withGoogleMap
        )(props =>
            <GoogleMap defaultZoom={16} center={props.center}>
                <Marker position={props.center} />
                <StreetViewPanorama position={props.center} visible />
            </GoogleMap>
    )

    handleExpandClick = () => {
        this.setState({ expanded: !this.state.expanded });
    };

    componentDidMount = async () => {
        await this.service.fetchHouseInfo(this.props.match.params['h_id']).then((result) => {this.setState({ info: result })});
        this.service.getHouseAddress(this.state.info.latitude, this.state.info.longitude).then((result) => {this.setState({ addr: result })});
        this.userService.fetchUserInfo(this.state.info.u_id).then((result) => {this.setState({ userinfo: result })});
    }

    render() {
        const { classes } = this.props;

        return (
            <div>
                <NavBar />
                <Tooltip id="tooltip-fab" title="Add to Compare">
                <Button variant="fab" mini color="secondary" aria-label="add" className={classes.button}>
                    <AddIcon />
                </Button>
                </Tooltip>
                <Tooltip id="tooltip-fab" title="Edit House">
                <Button variant="fab" mini aria-label="edit" className={classes.button}>
                    <EditIcon/>
                </Button>
                </Tooltip>
                <Tooltip id="tooltip-fab" title="Delete House">
                    <Button variant="fab" mini aria-label="delete" className={classes.button}>
                        <DeleteIcon />
                    </Button>
                </Tooltip>

                <Card className={classes.card}>
                    <CardHeader
                        avatar={
                            <Avatar className={classes.avatar}>
                                {this.state.userinfo.username.charAt(0)}
                            </Avatar>
                        }
                        action={
                            <CardActions>
                                <Button size="small">CONTACT SELLER</Button>
                                <Button variant='raised' color='primary' size="small" onClick={() => 
                                    this.service.buyHouse(this.state.info.h_id)
                                        .then(() => {
                                            alert('Congratulations!! The house is yours!!');
                                            this.props.history.replace('/home');
                                        })
                                        .catch((err) => alert('Something went wrong, please try again.'))}>
                                    BUY
                                </Button>
                            </CardActions>
                        }
                        title = {this.state.userinfo.username}
                        subheader = {this.state.userinfo.email}
                    />

                    <this.StreetViewPanorama center={ {lat: this.state.info.latitude, lng: this.state.info.longitude} } />
                    <CardContent>
                        <Typography variant="headline">
                            {this.state.addr}
                        </Typography>
                        <Typography color = "primary" variant="headline">
                        <Grid container spacing={8}>
                            <Grid item xs>
                            {this.state.info.bedroomCnt} BEDROOMS
                            </Grid>
                            <Grid item xs>
                            {this.state.info.bathroomCnt} BATHROOMS
                            </Grid>
                            <Grid item xs>
                            {this.state.info.livingAreaSize} SQFT
                            </Grid>
                        </Grid>
                        </Typography>
                        <Typography color = "secondary" variant="display1">
                            ${this.state.info.price}
                        </Typography>  
                    </CardContent>

                    <CardActions className={classes.actions} disableActionSpacing>
                        <IconButton aria-label="Add to favorites">
                        <FavoriteIcon />
                        </IconButton>
                        <IconButton aria-label="Viewed">
                        <VisibilityIcon />
                        </IconButton>
                        <IconButton
                        className={classnames(classes.expand, {
                            [classes.expandOpen]: this.state.expanded,
                        })}
                        onClick={this.handleExpandClick}
                        aria-expanded={this.state.expanded}
                        aria-label="Show more"
                        >
                        <ExpandMoreIcon />
                        </IconButton>
                    </CardActions>
                    <Collapse in={this.state.expanded} timeout="auto" unmountOnExit>
                        <CardContent>
                        <Typography paragraph variant="title">
                            Detailed Information:
                        </Typography>
                        <Typography paragraph>Year Built: {this.state.info.yearBuilt}</Typography>
                        <Typography>Building Quality: {this.state.info.buildingQualityID}</Typography>
                        <Typography paragraph variant="caption">(from 1 to 10 with 10 the best)</Typography>
                        <Typography paragraph>Number of Stories: {this.state.info.storyNum}</Typography>
                        <Typography paragraph>Lot Size: {this.state.info.lotSize} Sqft</Typography>
                        <Typography paragraph>Total Price: ${this.state.info.price}</Typography>
                        <Typography paragraph>Tax: ${this.state.info.tax} per year</Typography>
                        </CardContent>
                    </Collapse>
                </Card>

            </div>
        )
    }
}

HouseInfo.propTypes = {
    classes: PropTypes.object.isRequired,
};

export default withStyles(styles)(HouseInfo);
