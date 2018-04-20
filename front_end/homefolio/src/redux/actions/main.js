import { 
    ADD_COMPARE_HOUSE,
    EMPTY_COMPARE_HOUSES,
    HOME_MAP_BOUND_CHANGED,
    HOME_DISPLAY_HOUSES_CHANGED,
    HOME_MAP_FOCUSED_MARKER_CHANGED
} from "../constants/action-types";

export const addCompareHouses = house => ({ type: ADD_COMPARE_HOUSE, payload: house });
export const emptyCompareHouses = () => ({ type: EMPTY_COMPARE_HOUSES });
export const homeMapBoundChanged = newBounds => ({ type: HOME_MAP_BOUND_CHANGED, payload: newBounds });
export const homeDisplayHousesChanged = newHouses => ({ type: HOME_DISPLAY_HOUSES_CHANGED, payload: newHouses });
export const homeMapFocusedMarkerChanged = newFocused => ({ type: HOME_MAP_FOCUSED_MARKER_CHANGED, payload: newFocused });