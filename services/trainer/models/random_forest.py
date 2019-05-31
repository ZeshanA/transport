import logging

from sklearn.ensemble import RandomForestRegressor

from lib.models import SKModel


class RandomForestModel(SKModel):

    def __init__(self, route_id, custom_params):
        self.params = {'n_estimators': 200}
        super().__init__(route_id, custom_params)

    def __create_model__(self):
        logging.info("Creating model...")
        model = RandomForestRegressor(
            n_estimators=self.params['n_estimators'],
            n_jobs=-1
        )
        self.model = model
