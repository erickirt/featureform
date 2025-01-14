// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Copyright 2024 FeatureForm Inc.
//

import Container from '@mui/material/Container';
import Paper from '@mui/material/Paper';
import React, { useEffect } from 'react';
import Loader from 'react-loader-spinner';
import { connect } from 'react-redux';
import Resource from '../../api/resources/Resource.js';
import NotFound from '../../components/notfoundpage/NotFound.js';
import VariantNotFound from '../../components/notfoundpage/VariantNotFound.js';
import { setVariant } from '../resource-list/VariantSlice.js';
import { fetchEntity } from './EntityPageSlice.js';
import EntityPageView from './EntityPageView.js';

const mapDispatchToProps = (dispatch) => {
  return {
    fetch: (api, type, title) => dispatch(fetchEntity({ api, type, title })),
    setVariant: (type, name, variant) =>
      dispatch(setVariant({ type, name, variant })),
  };
};

function mapStateToProps(state) {
  return {
    entityPage: state.entityPage,
    activeVariants: state.selectedVariant,
  };
}

export function LoadingDots() {
  return (
    <div data-testid='loadingDotsId'>
      <Container maxWidth='xl'>
        <Paper elevation={3}>
          <Container style={{ textAlign: 'center' }} maxWidth='sm'>
            <Loader type='ThreeDots' color='grey' height={40} width={40} />
          </Container>
        </Paper>
      </Container>
    </div>
  );
}

const fetchNotFound = (object) => {
  return !object?.resources?.name && !object?.resources?.type;
};

const variantNotFound = (object, queryVariant = '') => {
  let result = true;
  if (queryVariant) {
    result = object?.resources?.['all-variants']?.includes(queryVariant);
  }
  return !result;
};

const EntityPage = ({
  api,
  entityPage,
  activeVariants,
  type,
  entity,
  queryVariant,
  ...props
}) => {
  let resourceType = Resource[Resource.pathToType[type]];
  const fetchEntity = props.fetch;

  useEffect(() => {
    const getData = async () => {
      if (api && type && entity) {
        await fetchEntity(api, type, entity);
      }
    };
    getData();
  }, [type, entity]);

  let body = <></>;
  if (entityPage.loading === true) {
    body = <LoadingDots />;
  } else if (entityPage.loading === false) {
    if (entityPage.failed === true || fetchNotFound(entityPage)) {
      body = <NotFound type={type} entity={entity} />;
    } else if (variantNotFound(entityPage, queryVariant)) {
      body = (
        <VariantNotFound
          type={type}
          entity={entity}
          queryVariant={queryVariant}
        />
      );
    } else {
      body = (
        <EntityPageView
          api={api}
          entity={entityPage}
          setVariant={props.setVariant}
          activeVariants={activeVariants}
          typePath={type}
          resourceType={resourceType}
          queryVariant={queryVariant}
        />
      );
    }
  }
  return body;
};

export default connect(mapStateToProps, mapDispatchToProps)(EntityPage);
