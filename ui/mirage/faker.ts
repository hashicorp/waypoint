/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import faker from 'faker';

// We want consistent fake data in development for easy page reloads etc.
faker.seed(1);

export default faker;
